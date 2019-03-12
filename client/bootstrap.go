package client

import (
	"github.com/alexhunt7/gofigure/utils"
	"github.com/alexhunt7/ssher"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"path"
	"strings"
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

func Bootstrap(host, configFile string) (*grpc.ClientConn, error) {
	var grpcConn *grpc.ClientConn

	sshConfig, connectString, err := ssher.ClientConfig(host, configFile)
	if err != nil {
		return grpcConn, err
	}

	sshConn, err := ssh.Dial("tcp", connectString, sshConfig)
	if err != nil {
		return grpcConn, err
	}
	defer sshConn.Close()

	// TODO kill existing process?

	sftpClient, err := sftp.NewClient(sshConn)
	if err != nil {
		return grpcConn, err
	}
	defer sftpClient.Close()

	// TODO pass this in?
	executable := path.Base(os.Args[0])
	err = putfile(sftpClient, os.Args[0], executable, 0700)
	if err != nil {
		return grpcConn, err
	}

	// TODO pass these in
	for _, filename := range []string{"ca-cert.pem", "cert.pem", "key.pem"} {
		err = putfile(sftpClient, "testdata/"+filename, filename, 0600)
		if err != nil {
			return grpcConn, err
		}
	}

	session, err := sshConn.NewSession()
	if err != nil {
		return grpcConn, err
	}
	defer session.Close()

	err = session.Start("./" + executable + " serve --caFile ca-cert.pem --certFile cert.pem --keyFile key.pem </dev/null >/dev/null 2>&1")
	if err != nil {
		return grpcConn, err
	}

	splitConnectString := strings.Split(connectString, ":")
	// TODO handle alternative ports
	grpcConn, err = ConnectGRPC(splitConnectString[0]+":10000", "testdata/ca-cert.pem", "testdata/cert.pem", "testdata/key.pem")
	if err != nil {
		return grpcConn, err
	}

	return grpcConn, nil
}

func ConnectGRPC(address, caFile, certFile, keyFile string) (*grpc.ClientConn, error) {
	creds, err := utils.LoadCredentials(caFile, certFile, keyFile)
	if err != nil {
		return nil, err
	}
	var conn *grpc.ClientConn
	for {
		log.Println("trying")
		conn, err = grpc.Dial(address, grpc.WithTransportCredentials(creds))
		if err == nil {
			break
		}
	}
	return conn, nil
}
