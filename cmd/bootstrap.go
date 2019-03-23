// Copyright © 2019 Alex Hunt <alex.s.hunt7@gmail.com>
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
	"context"
	"fmt"
	"github.com/alexhunt7/gofigure/master"
	pb "github.com/alexhunt7/gofigure/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/spf13/cobra"
	"log"
	"sync"
	"time"
)

// bootstrapCmd represents the bootstrap command
var (
	bootstrapCmd = &cobra.Command{
		Use:   "bootstrap [flags] host1 [host2] [host3] [...]",
		Args:  cobra.MinimumNArgs(1),
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("bootstrap called")

			clients, err := master.BootstrapMany(args, configFile)
			if err != nil {
				log.Fatal(err)
			}

			var wg sync.WaitGroup
			for _, client := range clients {
				client.RemoteDirectory("/home/alex/gofigure_dir")
				wg.Add(1000)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				for i := 0; i < 1000; i++ {
					go func(i int) {
						defer wg.Done()
						request := &pb.FileRequest{
							Properties: &pb.FileProperties{
								Path:  fmt.Sprintf("/home/alex/gofigure_dir/%d", i),
								Owner: "alex",
								Group: "alex",
								Mode:  "700",
							},
						}
						_, err := client.Directory(ctx, request, grpc_retry.WithMax(5))
						if err != nil {
							log.Printf("failed to create dir")
							log.Fatal(err)
						}
						log.Print(i)
					}(i)
				}
			}
			wg.Wait()
		},
	}
	configFile string
)

func init() {
	rootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().StringVar(&configFile, "config", "", "SSH config file")
}
