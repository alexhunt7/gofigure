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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"os"
)

var (
	caFile   string
	certFile string
	keyFile  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	//caFile := cmd.Flag("caFile").Value.String()
	//	//certFile := cmd.Flag("certFile").Value.String()
	//	//keyFile := cmd.Flag("keyFile").Value.String()
	//	fmt.Printf("CA %s\n", caFile)
	//	fmt.Printf("cert %s\n", certFile)
	//	fmt.Printf("key %s\n", keyFile)
	//},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&caFile, "caFile", "", "Path to CA certificate file.")
	rootCmd.PersistentFlags().StringVar(&certFile, "certFile", "", "Path to certificate file.")
	rootCmd.PersistentFlags().StringVar(&keyFile, "keyFile", "", "Path to certificate key file.")
	rootCmd.MarkPersistentFlagRequired("caFile")
	rootCmd.MarkPersistentFlagRequired("certFile")
	rootCmd.MarkPersistentFlagRequired("keyFile")
}

func loadCredentials(caFile string, certFile string, keyFile string) (credentials.TransportCredentials, error) {
	// https://github.com/grpc/grpc-go/issues/403

	peerCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Errorf("load peer cert/key error: %v", err)
		return nil, err
	}
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		fmt.Errorf("read ca cert file error: %v", err)
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	creds := credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{peerCert},
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})
	return creds, nil
}
