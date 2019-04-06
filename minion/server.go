/*
Copyright 2018 Alex Hunt

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Package minion implements the remote side of the gofigure
configuration management tool.

It includes functions for creating and deleting files and directories,
as well as process execution.
*/
package minion

import (
	"context"
	"fmt"
	"github.com/alexhunt7/gofigure/credentials"
	pb "github.com/alexhunt7/gofigure/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Minion implements the remote side of the gofigure service.
type Minion struct {
	grpcServer *grpc.Server
}

func (s *Minion) Exit(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	go s.grpcServer.GracefulStop()
	return &pb.Empty{}, nil
}

func Serve(caFile, certFile, keyFile string, bind net.IP, port int) {
	log.Println("Serving gofigure with:")
	log.Printf("  CA:   %s\n", caFile)
	log.Printf("  cert: %s\n", certFile)
	log.Printf("  key:  %s\n", keyFile)
	log.Printf("  bind: %s\n", bind)
	log.Printf("  port: %d\n", port)

	creds, err := credentials.Load(caFile, certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bind, port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterGofigureServer(grpcServer, &Minion{grpcServer: grpcServer})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
