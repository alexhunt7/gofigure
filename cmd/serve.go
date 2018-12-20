// Copyright © 2018 Alex Hunt <alex.s.hunt7@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	pb "github.com/alexhunt7/gofigure/proto"
	server "github.com/alexhunt7/gofigure/server"
	"log"
	"net"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// serveCmd represents the serve command
var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Listen for client connections.",
		Long: `A longer description that spans multiple lines and likely contains examples
and usag	e of using your command. For example:

Cobra is	 a CLI library for Go that empowers applications.
This app	lication is a tool to generate the needed files
to quick	ly create a Cobra application.`,
		Run: serve,
	}
	bind string
	port int
)

func serve(cmd *cobra.Command, args []string) {
	fmt.Println("Serving gofigure with:")
	fmt.Printf("  CA:   %s\n", caFile)
	fmt.Printf("  cert: %s\n", certFile)
	fmt.Printf("  key:  %s\n", keyFile)
	fmt.Printf("  bind: %s\n", bind)
	fmt.Printf("  port: %d\n", port)

	creds, err := loadCredentials(caFile, certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bind, port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterGofigureServer(grpcServer, &server.GofigureServer{})
	grpcServer.Serve(lis)
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&bind, "bind", "127.0.0.1", "Address to listen on.")
	serveCmd.Flags().IntVar(&port, "port", 10000, "TCP Port to listen on.")
}
