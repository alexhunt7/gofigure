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

// Package master is the gofigure reference client library with helper functions.
package master

import (
	pb "github.com/alexhunt7/gofigure/proto"
)

// Client is a concrete example implementation of the local side of the gofigure
// configuration management tool.
type Client struct {
	pb.GofigureClient
}
