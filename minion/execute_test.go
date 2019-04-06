package minion

import (
	"context"
	"testing"

	pb "github.com/alexhunt7/gofigure/proto"
)

//https://stackoverflow.com/questions/42102496/testing-a-grpc-service
func TestExec(t *testing.T) {
	minion := Minion{}

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
			// Simple case
			executable:         "echo",
			args:               []string{"hello", "world"},
			expectedReturnCode: 0,
			expectedStdout:     "hello world\n",
			expectedStderr:     "",
		},
		{
			// test it works with no args
			executable:         "echo",
			expectedReturnCode: 0,
			expectedStdout:     "\n",
			expectedStderr:     "",
		},
		{
			// test it gives stderr, passes return code, and can handle shells
			executable:         "sh",
			args:               []string{"-c", "echo hello world && echo goodbye world >&2 && exit 37"},
			expectedReturnCode: 37,
			expectedStdout:     "hello world\n",
			expectedStderr:     "goodbye world\n",
		},
		{
			// test it throws an err with a non-existent executable
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
		resp, err := minion.Exec(context.Background(), &req)
		if tt.shouldErr {
			if err == nil {
				t.Errorf("TestExec did not error as expected")
			}
			continue
		} else {
			if err != nil {
				t.Errorf("TestExec got unexpected error: %v", err)
			}
		}
		if resp.ReturnCode != tt.expectedReturnCode {
			t.Errorf("TestExec return code, %v != %v", tt.expectedReturnCode, resp.ReturnCode)
		}
		if resp.Stdout != tt.expectedStdout {
			t.Errorf("TestExec stdout, %v != %v", tt.expectedStdout, resp.Stdout)
		}
		if resp.Stderr != tt.expectedStderr {
			t.Errorf("TestExec stderr, %v != %v", tt.expectedStderr, resp.Stderr)
		}
	}
}
