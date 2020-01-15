package master

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/ghodss/yaml"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alexhunt7/gofigure/credentials"
	pb "github.com/alexhunt7/gofigure/proto"
)

var (
	Debug = true
)

type Creds struct {
	CAFile   string
	CertFile string
	KeyFile  string
}

type MinionConfig struct {
	Bind  net.IP
	Port  int
	Creds *Creds
}

type Config struct {
	Defaults *MinionConfig
	Master   *Creds
	Minions  map[string]*MinionConfig
}

func (c *Config) promoteDefaults() {
	if c.Defaults.Creds == nil {
		c.Defaults.Creds = &Creds{}
	}
	for _, minionConfig := range c.Minions {
		if minionConfig.Bind == nil {
			minionConfig.Bind = c.Defaults.Bind
		}
		if minionConfig.Port == 0 {
			minionConfig.Port = c.Defaults.Port
		}
		if minionConfig.Creds == nil {
			minionConfig.Creds = &Creds{}
		}
		if minionConfig.Creds.CAFile == "" {
			minionConfig.Creds.CAFile = c.Defaults.Creds.CAFile
		}
		if minionConfig.Creds.CertFile == "" {
			minionConfig.Creds.CertFile = c.Defaults.Creds.CertFile
		}
		if minionConfig.Creds.KeyFile == "" {
			minionConfig.Creds.KeyFile = c.Defaults.Creds.KeyFile
		}
	}
}

func (c *Config) parse(filename string) error {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}

	err = yaml.Unmarshal(f, c)
	if err != nil {
		return fmt.Errorf("error unmarshalling yaml: %v", err)
	}
	c.promoteDefaults()
	// TODO confirm nothing is nil
	return nil
}

// Bootstrap will parse an openssh config file, ssh to the remote host, copy the executable there,
// run it, and attempt to connect, returning a gofigure client.
func Bootstrap(host, sshConfigPath, executable string, minionConfig *MinionConfig, masterCreds *Creds) (*Client, error) {
	// TODO remove existing servers
	var gofigureClient *Client
	args := []string{
		executable,
		minionConfig.Creds.CAFile,
		minionConfig.Creds.CertFile,
		minionConfig.Creds.KeyFile,
		fmt.Sprintf("%s:", host),
	}
	if sshConfigPath != "" {
		args = append([]string{"-F", sshConfigPath}, args...)
	}

	cmd := exec.Command("scp", args...)
	if Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		return gofigureClient, fmt.Errorf("failed to scp minion executable and certificate files: %v", err)
	}

	remoteExec := fmt.Sprintf("./%s", path.Base(executable))
	// TODO allow arg for sudo
	args = []string{
		host,
		fmt.Sprintf("sudo %s --bind %s --port %d --caFile %s --certFile %s --keyFile %s",
			remoteExec,
			minionConfig.Bind.String(),
			minionConfig.Port,
			path.Base(minionConfig.Creds.CAFile),
			path.Base(minionConfig.Creds.CertFile),
			path.Base(minionConfig.Creds.KeyFile)),
	}
	if sshConfigPath != "" {
		args = append([]string{"-F", sshConfigPath}, args...)
	}
	cmd = exec.Command("ssh", args...)
	if Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err = cmd.Start()
	if err != nil {
		return gofigureClient, fmt.Errorf("failed to execute %s on remote host %s: %v", remoteExec, host, err)
	}

	grpcConnectString := fmt.Sprintf("%s:%d", host, minionConfig.Port)

	conn, err := ConnectGRPC(grpcConnectString, masterCreds.CAFile, masterCreds.CertFile, masterCreds.KeyFile)
	if err != nil {
		return gofigureClient, err
	}

	return &Client{GofigureClient: pb.NewGofigureClient(conn)}, nil
}

// ConnectGRPC attempts to connect over GRPC to the remote address.
// It returns only a GRPC connection, not a gofigure client.
func ConnectGRPC(address, caFile, certFile, keyFile string) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	tries := 1
	maxTries := 30
	//for i := 0; i < 30; i++ {
	for {
		c, err := net.Dial("tcp", address)
		if err == nil {
			// TODO reuse this connection instead of closing it
			c.Close()
			break
		}
		tries++
		if tries > maxTries {
			return nil, err
		}
		time.Sleep(100 * time.Millisecond)
	}

	creds, err := credentials.Load(caFile, certFile, keyFile)
	if err != nil {
		return nil, err
	}
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

// BootstrapMany calls Bootstrap for many hosts in parallel.
func BootstrapMany(sshConfigPath string, executable string, minionConfigs map[string]*MinionConfig, masterCreds *Creds) (map[string]*Client, error) {
	type result struct {
		host   string
		client *Client
	}
	gobootstrap := func(host, sshConfigPath, executable string, minionConfig *MinionConfig, masterCreds *Creds, successChan chan<- *result, failChan chan<- error) {
		client, err := Bootstrap(host, sshConfigPath, executable, minionConfig, masterCreds)
		if err != nil {
			failChan <- err
			return
		}
		successChan <- &result{host: host, client: client}
	}

	successChan, failChan := make(chan *result), make(chan error)
	for host, minionConfig := range minionConfigs {
		go gobootstrap(host, sshConfigPath, executable, minionConfig, masterCreds, successChan, failChan)
	}

	clients := make(map[string]*Client)
	for range minionConfigs {
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

func BootstrapCmd() (map[string]*Client, error) {
	clients := make(map[string]*Client)

	app := kingpin.New("name", "help")
	bootstrapConfig := app.Arg("config", "Config file.").Required().String()
	minionPath := app.Arg("minion", "Path to minion binary").Required().String()
	_, err := app.Parse(os.Args[1:])
	if err != nil {
		return clients, err
	}

	var config Config
	err = config.parse(*bootstrapConfig)
	if err != nil {
		return clients, err
	}

	clients, err = BootstrapMany("", *minionPath, config.Minions, config.Master)
	if err != nil {
		return clients, err
	}
	return clients, nil
}

func Exit(clients map[string]*Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return RunAll(ctx, clients, func(client *Client) error {
		_, err := client.Exit(ctx, &pb.Empty{}, grpc_retry.WithMax(5))
		return err
	})
}

func DeferredExit(clients map[string]*Client) {
	err := Exit(clients)
	if err != nil {
		log.Fatalf("failed to shutdown gofigure minions: %v", err)
	}
}
