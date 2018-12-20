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
	gclient "github.com/alexhunt7/gofigure/client"
	"google.golang.org/grpc"
	"github.com/spf13/cobra"
	"log"
)

// dialCmd represents the dial command
var (
	dialCmd = &cobra.Command{
		Use:   "dial",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usag	e of using your command. For example:

Cobra is	 a CLI library for Go that empowers applications.
This app	lication is a tool to generate the needed files
to quick	ly create a Cobra application.`,
		Run: dial,
	}
	address string
)

func dial(cmd *cobra.Command, args []string) {
	fmt.Println("dial called")
	creds, err := loadCredentials(caFile, certFile, keyFile)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
        if err != nil {
                log.Fatalf("fail to dial: %v", err)
        }
        defer conn.Close()
	client := &gclient.Client{pb.NewGofigureClient(conn)}
	client.CreateDir("/home/alex/go/src/github.com/alexhunt7/gofigure/asdf")
	client.CreateFile("/home/alex/go/src/github.com/alexhunt7/gofigure/asdf/qwer")
	client.Stat("/home/alex/go/src/github.com/alexhunt7/gofigure/asdf/qwer")
	client.Exec("echo", "hello", "world")
}

func init() {
	rootCmd.AddCommand(dialCmd)
	dialCmd.Flags().StringVar(&address, "address", "127.0.0.1:10000", "Address of gofigure server.")
}
