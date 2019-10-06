package minion

import (
	"fmt"
	"net"
	"testing"
	"github.com/stretchr/testify/require"

)


func TestMetadata(t *testing.T) {
	require := require.New(t)

	caFile := "../testdata/ca-cert.pem"
	certFile := "../testdata/cert.pem"
	keyFile := "../testdata/key.pem"
	bind := net.ParseIP("127.0.0.1")
	port := 10001

	testClient, err := ConstructMinion(nil,
		nil,
		CreateGrpcServer(caFile, certFile, keyFile),
		SetBind(bind),
		SetPort(port),
	)
	if testClient.Metadata.Port != port {
		t.Fatal(fmt.Sprintf("%d != %d", port, testClient.Metadata.Port))
	}

	if testClient.Metadata.CertFile != certFile{
		t.Fatal(fmt.Sprintf("%s != %s", certFile, testClient.Metadata.CertFile))
	}

	if testClient.Metadata.CaFile != caFile{
		t.Fatal(fmt.Sprintf("%s != %s", caFile, testClient.Metadata.CaFile))
	}

	if testClient.Metadata.KeyFile != keyFile{
		t.Fatal(fmt.Sprintf("%s != %s", keyFile, testClient.Metadata.KeyFile))
	}

	require.Nil(err)
}