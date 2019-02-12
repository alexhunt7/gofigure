package client

import (
	"github.com/alexhunt7/ssher"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
)

func Bootstrap(host, configFile string, successChan chan<- ssh.Conn, failChan chan<- error) {
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
	w, err := sftpClient.Create(executable)
	if err != nil {
		failChan <- err
		return
	}

	r, err := os.Open(os.Args[0])
	if err != nil {
		failChan <- err
		return
	}

	_, err = io.Copy(w, r)
	if err != nil {
		failChan <- err
		return
	}

	// TODO fsync?

	err = sftpClient.Chmod(executable, 0700)
	if err != nil {
		failChan <- err
		return
	}

	successChan <- conn
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
