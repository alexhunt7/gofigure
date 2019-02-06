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
	"github.com/alexhunt7/ssher"
	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"path"
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
			log.Println(len(args))
			connectSuccess, connectFailure := make(chan int), make(chan error)
			for i, host := range args {
				log.Println(i, host)
				go func(host, configFile string, i int) {
					conn, err := connect(host, configFile)
					if err != nil {
						connectFailure <- err
						return
					}
					defer conn.Close()

					sftpClient, err := sftp.NewClient(conn)
					if err != nil {
						connectFailure <- err
						return
					}
					defer sftpClient.Close()

					executable := path.Base(os.Args[0])
					w, err := sftpClient.Create(executable)
					if err != nil {
						connectFailure <- err
						return
					}

					r, err := os.Open(os.Args[0])
					if err != nil {
						connectFailure <- err
						return
					}

					_, err = io.Copy(w, r)
					if err != nil {
						connectFailure <- err
						return
					}

					connectSuccess <- i
				}(host, configFile, i)
			}
			for i := 0; i < len(args); i++ {
				select {
				case res := <-connectSuccess:
					log.Println(res)
				case err := <-connectFailure:
					log.Println(err)
				}
			}
		},
	}
	configFile string
)

func connect(host, configFile string) (*ssh.Client, error) {
	config, connectString, err := ssher.ClientConfig(host, configFile)
	if err != nil {
		return nil, err
	}
	conn, err := ssh.Dial("tcp", connectString, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().StringVar(&configFile, "config", "", "SSH config file")
}
