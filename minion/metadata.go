package minion

import (
	"github.com/alexhunt7/gofigure/credentials"
	"google.golang.org/grpc"
	"net"
	"sync"
)

// Minion implements the remote side of the gofigure service.
// It also is a threadsafe struct holding all meta info
type Minion struct {
	Lock       *sync.Mutex
	grpcServer *grpc.Server
	Metadata   *MinionMetadata
}

type MinionMetadata struct {
	CaFile   string
	CertFile string
	KeyFile  string
	Bind     net.IP
	Port     int
}

func (minion *MinionMetadata) Enumerate() {}

// Constructs the minion right away if we have the grpcServer variable + metadata
func ConstructMinion(lock *sync.Mutex, grpcServer *grpc.Server, optionalVars ...func(*Minion) error) (*Minion, error) {
	obj := &Minion{
		Lock:       lock,
		grpcServer: grpcServer,
		Metadata:   &MinionMetadata{},
	}

	// Create mutex lock if we're not passing in a lock already
	if obj.Lock == nil {
		obj.Lock = &sync.Mutex{}
	}

	// for any other optional options the user can pass into the constructor
	for _, setVars := range optionalVars {
		err := setVars(obj)

		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

func SetCaFile(caFile string) func(*Minion) error {
	return func(minion *Minion) error {
		if caFile != "" {
			minion.Lock.Lock()
			minion.Metadata.CaFile = caFile
			minion.Lock.Unlock()
		}
		return nil
	}
}

func SetCertFile(certFile string) func(*Minion) error {
	return func(minion *Minion) error {
		if certFile != "" {
			minion.Lock.Lock()
			minion.Metadata.CertFile = certFile
			minion.Lock.Unlock()
		}
		return nil
	}
}
func SetKeyFile(keyFile string) func(*Minion) error {
	return func(minion *Minion) error {
		if keyFile != "" {
			minion.Lock.Lock()
			minion.Metadata.KeyFile = keyFile
			minion.Lock.Unlock()
		}
		return nil
	}
}
func SetBind(bind net.IP) func(*Minion) error {
	return func(minion *Minion) error {
		if bind != nil {
			minion.Lock.Lock()
			minion.Metadata.Bind = bind
			minion.Lock.Unlock()
		}
		return nil
	}
}
func SetPort(port int) func(*Minion) error {
	return func(minion *Minion) error {
		if port != 0 {
			minion.Lock.Lock()
			minion.Metadata.Port = port
			minion.Lock.Unlock()
		}
		return nil
	}
}

// Function to set the entire object if the user passes in all metadata struct
func CreateGrpcServer(caFile, certFile, keyFile string) func(*Minion) error {
	return func(minion *Minion) error {
		creds, err := credentials.Load(caFile, certFile, keyFile)
		if err != nil {
			return err
		}

		grpcServer := grpc.NewServer(grpc.Creds(creds))

		minion.Lock.Lock()
		minion.grpcServer = grpcServer
		minion.Lock.Unlock()

		minion.Lock.Lock()
		minion.Metadata.CaFile = caFile
		minion.Lock.Unlock()

		minion.Lock.Lock()
		minion.Metadata.KeyFile = keyFile
		minion.Lock.Unlock()

		minion.Lock.Lock()
		minion.Metadata.CertFile = certFile
		minion.Lock.Unlock()

		return nil
	}
}
