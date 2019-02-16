package client

import (
	"github.com/alexhunt7/ssher"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
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

func Bootstrap(host, configFile string, successChan chan<- string, failChan chan<- error) {
	conn, err := Connect(host, configFile)
	if err != nil {
		failChan <- err
		return
	}
	defer conn.Close()

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		failChan <- err
		return
	}
	defer sftpClient.Close()

	// TODO pass this in?
	executable := path.Base(os.Args[0])
	err = putfile(sftpClient, os.Args[0], executable, 0700)
	if err != nil {
		failChan <- err
		return
	}

	// TODO pass these in
	for _, filename := range []string{"ca-cert.pem", "cert.pem", "key.pem"} {
		err = putfile(sftpClient, "testdata/"+filename, filename, 0600)
		if err != nil {
			failChan <- err
			return
		}
	}

	session, err := conn.NewSession()
	if err != nil {
		failChan <- err
		return
	}
	defer session.Close()

	err = session.Start("./" + executable + " serve --caFile ca-cert.pem --certFile cert.pem --keyFile key.pem")
	if err != nil {
		failChan <- err
		return
	}

	time.Sleep(1 * time.Second)

	successChan <- host
}

func Connect(host, configFile string) (*ssh.Client, error) {
	config, connectString, err := ssher.ClientConfig(host, configFile)
	if err != nil {
		return nil, err
	}
	conn, err := ssh.Dial("tcp", connectString, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
