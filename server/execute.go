package main

import (
	"bytes"
	pb "github.com/alexhunt7/gofigure/proto"
	"golang.org/x/net/context"
	"os/exec"
	"syscall"
)

func (s *GofigureServer) GofigureExec(ctx context.Context, req *pb.ExecRequest) (*pb.ExecResult, error) {
	// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
	// TODO env
	// TODO stdin
	// TODO cwd
	cmd := exec.CommandContext(ctx, req.Executable, req.Args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	rc := 255
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
			Stdout: string(stdout.Bytes()),
			Stderr: string(stderr.Bytes())},
		nil
}
