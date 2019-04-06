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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	//"sync"
	"time"

	"github.com/ghodss/yaml"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alexhunt7/gofigure/master"
	"github.com/alexhunt7/gofigure/minion"
	pb "github.com/alexhunt7/gofigure/proto"
)

var (
	bootstrap       = kingpin.Command("bootstrap", "Deploy minions on remote machines, and connect to them.")
	bootstrapConfig = bootstrap.Arg("config", "Config file.").Required().String()

	serve         = kingpin.Command("serve", "Listen for remote commands as a minion.")
	serveCAFile   = serve.Flag("caFile", "Path to CA file.").Required().String()
	serveCertFile = serve.Flag("certFile", "Path to certificate file.").Required().String()
	serveKeyFile  = serve.Flag("keyFile", "Path to key file.").Required().String()
	serveBind     = serve.Flag("bind", "Address to listen on.").Required().IP()
	servePort     = serve.Flag("port", "Port to listen on.").Required().Int()
)

type Config struct {
	Defaults *master.MinionConfig
	Master   *master.Creds
	Minions  map[string]*master.MinionConfig
}

func (c *Config) promoteDefaults() {
	if c.Defaults.Creds == nil {
		c.Defaults.Creds = &master.Creds{}
	}
	for _, minionConfig := range c.Minions {
		if minionConfig.Bind == nil {
			minionConfig.Bind = c.Defaults.Bind
		}
		if minionConfig.Port == 0 {
			minionConfig.Port = c.Defaults.Port
		}
		if minionConfig.Creds == nil {
			minionConfig.Creds = &master.Creds{}
		}
		if minionConfig.Creds.CAFile == "" {
			minionConfig.Creds.CAFile = c.Defaults.Creds.CAFile
		}
		if minionConfig.Creds.CertFile == "" {
			minionConfig.Creds.CertFile = c.Defaults.Creds.CertFile
		}
		if minionConfig.Creds.KeyFile == "" {
			minionConfig.Creds.KeyFile = c.Defaults.Creds.KeyFile
		}
	}
}

func (c *Config) parse(filename string) error {
	f, err := ioutil.ReadFile("testdata/config.yml")
	if err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}

	err = yaml.Unmarshal(f, c)
	if err != nil {
		return fmt.Errorf("error unmarshalling yaml: %v", err)
	}
	c.promoteDefaults()
	// TODO confirm nothing is nil
	return nil
}

func main() {
	log.SetOutput(os.Stdout)
	switch kingpin.Parse() {
	case "bootstrap":
		var config Config
		err := config.parse(*bootstrapConfig)
		if err != nil {
			log.Fatal(err)
		}
		clients, err := master.BootstrapMany("", os.Args[0], config.Minions, config.Master)
		if err != nil {
			log.Fatal(err)
		}
		err = createDirs(clients)
		if err != nil {
			log.Fatal(err)
		}
	case "serve":
		minion.Serve(*serveCAFile, *serveCertFile, *serveKeyFile, *serveBind, *servePort)
	}
}

func createDirs(clients map[string]*master.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := master.RunAll(ctx, clients, func(client *master.Client) error {
		request := &pb.FileRequest{
			Properties: &pb.FileProperties{
				Path:  "/home/alex/gofigure_dir",
				Owner: "alex",
				Group: "alex",
				Mode:  "700",
			},
		}
		_, err := client.Directory(ctx, request, grpc_retry.WithMax(5))
		return err
	})
	if err != nil {
		return err
	}
	err = master.RunAll(ctx, clients, func(client *master.Client) error {
		for i := 0; i < 1000; i++ {
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
				return err
			}
		}
		return nil
	})
	return err
}
