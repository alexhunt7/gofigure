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
	"github.com/alexhunt7/gofigure/client"
	"github.com/spf13/cobra"
	"log"
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
			successChan, failChan := make(chan string), make(chan error)
			for i, host := range args {
				log.Println(i, host)
				go client.Bootstrap(host, configFile, successChan, failChan)
			}
			for i := 0; i < len(args); i++ {
				select {
				case res := <-successChan:
					log.Println("success " + res)
				case err := <-failChan:
					log.Println(err)
				}
			}
		},
	}
	configFile string
)

func init() {
	rootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().StringVar(&configFile, "config", "", "SSH config file")
}
