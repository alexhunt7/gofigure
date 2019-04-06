package master

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBootstrap(t *testing.T) {
	port := "32080"

	testdataDir, err := filepath.Abs("../testdata")
	if err != nil {
		t.Error(err)
		return
	}
	knownHosts := testdataDir + "/docker_known_hosts"
	identity := testdataDir + "/ssh_host_rsa_key"
	sshConfigPath := testdataDir + "/docker"

	containerIDBytes, err := exec.Command("docker", "run", "-d", "--rm",
		"-p", port+":22",
		"-p", "10000:10000",
		"gofigure-sshd").Output()
	if err != nil {
		t.Error(err)
		return
	}
	containerID := strings.TrimSpace(string(containerIDBytes))

	hasConnected := false
	cleanup := func() {
		// Docker waits an extra 10 seconds if it hasn't finished launching
		// the process when we try to stop it.
		if !hasConnected {
			time.Sleep(time.Millisecond * 100)
		}
		// TODO how to handle errors in a defer?
		_ = exec.Command("docker", "stop", containerID).Run()
	}
	defer cleanup()

	f, err := os.Create(sshConfigPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	_, err = f.WriteString("Host docker\nUser testuser\n    Hostname 127.0.0.1\n    Port " + port + "\n    UserKnownHostsFile " + knownHosts + "\n    IdentityFile " + identity + "\n")
	if err != nil {
		t.Error(err)
		return
	}
	err = f.Sync()
	if err != nil {
		t.Error(err)
		return
	}

	minionConfig := &MinionConfig{
		Bind: net.ParseIP("0.0.0.0"),
		Port: 10000,
		Creds: &Creds{
			CAFile:   fmt.Sprintf("%s/%s", testdataDir, "ca-cert.pem"),
			CertFile: fmt.Sprintf("%s/%s", testdataDir, "cert.pem"),
			KeyFile:  fmt.Sprintf("%s/%s", testdataDir, "key.pem"),
		},
	}

	masterCreds := &Creds{
		CAFile:   fmt.Sprintf("%s/%s", testdataDir, "ca-cert.pem"),
		CertFile: fmt.Sprintf("%s/%s", testdataDir, "cert.pem"),
		KeyFile:  fmt.Sprintf("%s/%s", testdataDir, "key.pem"),
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("failed to determine executable: %v", err)
	}

	_, err = BootstrapMany(sshConfigPath, executable, map[string]*MinionConfig{"docker": minionConfig}, masterCreds)
	if err != nil {
		t.Fatalf("failed to bootstrap: %v", err)
	}
}
