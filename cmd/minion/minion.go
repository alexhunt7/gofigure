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
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alexhunt7/gofigure/minion"
)

func main() {
	caFile := kingpin.Flag("caFile", "Path to CA file.").Required().String()
	certFile := kingpin.Flag("certFile", "Path to certificate file.").Required().String()
	keyFile := kingpin.Flag("keyFile", "Path to key file.").Required().String()
	bind := kingpin.Flag("bind", "Address to listen on.").Required().IP()
	port := kingpin.Flag("port", "Port to listen on.").Required().Int()

	kingpin.Parse()

	minion.Serve(*caFile, *certFile, *keyFile, *bind, *port)
}
