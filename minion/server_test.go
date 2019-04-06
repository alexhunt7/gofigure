package minion

import (
	"context"
	"net"
	"testing"
	"time"

	//"github.com/alexhunt7/gofigure/credentials"
	"github.com/alexhunt7/gofigure/master"
	pb "github.com/alexhunt7/gofigure/proto"
	"github.com/stretchr/testify/require"
)

func TestServe(t *testing.T) {
	require := require.New(t)

	caFile := "../testdata/ca-cert.pem"
	certFile := "../testdata/cert.pem"
	keyFile := "../testdata/key.pem"
	go Serve(caFile, certFile, keyFile, net.ParseIP("127.0.0.1"), 10001)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := master.ConnectGRPC("127.0.0.1:10001", caFile, certFile, keyFile)
	require.Nil(err)

	client := &master.Client{GofigureClient: pb.NewGofigureClient(conn)}
	_, err = client.Exit(ctx, &pb.Empty{})
	require.Nil(err)
}
