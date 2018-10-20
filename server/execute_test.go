package main

import (
	pb "github.com/alexhunt7/gofigure/proto"
	"golang.org/x/net/context"
	"testing"
)

//https://stackoverflow.com/questions/42102496/testing-a-grpc-service
func TestGofigureExec(t *testing.T) {
	s := GofigureServer{}

	tests := []struct {
		executable         string
		args               []string
		stdin              []byte
		env                []string
		shouldErr          bool
		expectedReturnCode int32
		expectedStdout     string
		expectedStderr     string
	}{
		{
			executable:         "echo",
			args:               []string{"hello", "world"},
			expectedReturnCode: 0,
			expectedStdout:     "hello world\n",
			expectedStderr:     "",
		},
		{
			executable: "asdflkjas",
			shouldErr:  true,
		},
	}

	for _, tt := range tests {
		req := pb.ExecRequest{
			Executable: tt.executable,
			Args:       tt.args,
			Stdin:      tt.stdin,
			Env:        tt.env,
		}
		resp, err := s.GofigureExec(context.Background(), &req)
		if tt.shouldErr {
			if err == nil {
				t.Errorf("TestGofigureExec did not error as expected")
			}
			continue
		} else {
			if err != nil {
				t.Errorf("TestGofigureExec got unexpected error: %v", err)
			}
		}
		if resp.ReturnCode != tt.expectedReturnCode {
			t.Errorf("TestGofigureExec return code, %v != %v", tt.expectedReturnCode, resp.ReturnCode)
		}
		if resp.Stdout != tt.expectedStdout {
			t.Errorf("TestGofigureExec stdout, %v != %v", tt.expectedStdout, resp.Stdout)
		}
		if resp.Stderr != tt.expectedStderr {
			t.Errorf("TestGofigureExec stderr, %v != %v", tt.expectedStderr, resp.Stderr)
		}
	}
}
