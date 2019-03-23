package master

import (
	pb "github.com/alexhunt7/gofigure/proto"
	"github.com/alexhunt7/gofigure/utils"
	"github.com/alexhunt7/ssher"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

func putfile(sftpClient *sftp.Client, src, dst string, perms os.FileMode) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := sftpClient.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	err = sftpClient.Chmod(dst, perms)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	// TODO fsync?
	return nil
}

func Bootstrap(host, configFile string) (*Client, error) {
	var gofigureClient *Client

	sshConfig, connectString, err := ssher.ClientConfig(host, configFile)
	if err != nil {
		return gofigureClient, err
	}

	sshConn, err := ssh.Dial("tcp", connectString, sshConfig)
	if err != nil {
		return gofigureClient, err
	}
	defer sshConn.Close()

	// TODO kill existing process?

	sftpClient, err := sftp.NewClient(sshConn)
	if err != nil {
		return gofigureClient, err
	}
	defer sftpClient.Close()

	// TODO pass this in?
	executable := path.Base(os.Args[0])
	err = putfile(sftpClient, os.Args[0], executable, 0700)
	if err != nil {
		return gofigureClient, err
	}

	// TODO pass these in
	for _, filename := range []string{"ca-cert.pem", "cert.pem", "key.pem"} {
		err = putfile(sftpClient, "testdata/"+filename, filename, 0600)
		if err != nil {
			return gofigureClient, err
		}
	}

	session, err := sshConn.NewSession()
	if err != nil {
		return gofigureClient, err
	}
	defer session.Close()

	// TODO wait for this to start?
	err = session.Start("./" + executable + " serve --caFile ca-cert.pem --certFile cert.pem --keyFile key.pem </dev/null >/dev/null 2>&1")
	if err != nil {
		return gofigureClient, err
	}

	splitConnectString := strings.Split(connectString, ":")
	// TODO handle alternative ports
	grpcConnectString := splitConnectString[0] + ":10000"

	tries := 1
	maxTries := 30
	//for i := 0; i < 30; i++ {
	for {
		c, err := net.Dial("tcp", grpcConnectString)
		if err == nil {
			// TODO reuse this connection instead of closing it
			c.Close()
			break
		}
		tries++
		if tries > maxTries {
			return gofigureClient, err
		}
		time.Sleep(100 * time.Millisecond)
	}

	conn, err := ConnectGRPC(grpcConnectString, "testdata/ca-cert.pem", "testdata/cert.pem", "testdata/key.pem")
	if err != nil {
		return gofigureClient, err
	}

	return &Client{GofigureClient: pb.NewGofigureClient(conn)}, nil
}

func ConnectGRPC(address, caFile, certFile, keyFile string) (*grpc.ClientConn, error) {
	creds, err := utils.LoadCredentials(caFile, certFile, keyFile)
	if err != nil {
		return nil, err
	}
	var conn *grpc.ClientConn
	for {
		conn, err = grpc.Dial(address,
			grpc.WithTransportCredentials(creds),
			grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor()),
			grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor()),
		)
		if err == nil {
			break
		}
	}
	return conn, nil
}

func BootstrapMany(hosts []string, configFile string) (map[string]*Client, error) {
	type result struct {
		host   string
		client *Client
	}

	successChan, failChan := make(chan *result), make(chan error)
	for _, host := range hosts {
		go func(host, configFile string, successChan chan<- *result, failChan chan<- error) {
			client, err := Bootstrap(host, configFile)
			if err != nil {
				failChan <- err
				return
			}
			successChan <- &result{host: host, client: client}
		}(host, configFile, successChan, failChan)
	}

	clients := make(map[string]*Client)
	for range hosts {
		select {
		case result := <-successChan:
			clients[result.host] = result.client
		case err := <-failChan:
			// TODO return multiple errors?
			return nil, err
		}
	}
	return clients, nil
}
