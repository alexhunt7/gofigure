package main

import (
	pb "alex/gofigure/proto"
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func runCreateDir(client pb.GofigureClient, path string) {
	log.Printf("runCreateDir")
	request := &pb.FileRequest{
		Properties: &pb.FileProperties{
			Path:  path,
			Owner: "alex",
			Group: "alex",
			Mode:  "666",
		},
	}
	log.Printf("Instantiated request")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("Created context")
	//response, err := client.GofigureDirectory(ctx, request)
	_, err := client.GofigureDirectory(ctx, request)
	log.Printf("ran client.GofigureDirectory")
	//log.Printf(response.Msg)
	if err != nil {
		log.Printf("failed to create dir")
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		//if *caFile == "" {
		//	*caFile = testdata.Path("ca.pem")
		//}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewGofigureClient(conn)

	runCreateDir(client, "/home/alex/git/golang/src/alex/gofigure/asdf")
}
