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

package minion

import (
	"bytes"
	"context"
	"os/exec"
	"syscall"

	pb "github.com/alexhunt7/gofigure/proto"
)

// Exec executes a command specified in an ExecRequest
// and returns the resulting ReturnCode, Stdout, and Stderr
// in an ExecResult.
func (s *Minion) Exec(ctx context.Context, req *pb.ExecRequest) (*pb.ExecResult, error) {
	// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
	// TODO env
	// TODO stdin
	// TODO cwd
	cmd := exec.CommandContext(ctx, req.Executable, req.Args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	var rc int
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// rc != 0
			ws := exiterr.Sys().(syscall.WaitStatus)
			rc = ws.ExitStatus()
		} else {
			return nil, err
		}
	} else {
		// should be rc == 0
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		rc = ws.ExitStatus()
	}

	return &pb.ExecResult{ReturnCode: int32(rc),
			Stdout: stdout.String(),
			Stderr: stderr.String()},
		nil
}
