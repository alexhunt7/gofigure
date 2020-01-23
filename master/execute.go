package master

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/alexhunt7/gofigure/proto"
)

func CheckExec(client *Client, args []string, timeoutSecs int) (*pb.ExecResult, error) {
	res, err := Exec(client, args, timeoutSecs)
	if err == nil && int(res.ReturnCode) != 0 {
		return res, fmt.Errorf("non-zero return code: %d\n\nSTDOUT:\n%s\n\nSTDERR:\n%s", res.ReturnCode, res.Stdout, res.Stderr)
	}
	return res, err
}

// You probably want CheckExec instead. This one won't error for non-zero return codes
func Exec(client *Client, args []string, timeoutSecs int) (*pb.ExecResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	request := &pb.ExecRequest{
		Executable: args[0],
		Args:       args[1:],
	}

	return client.Exec(ctx, request)
}

func Lines(s string) []string {
	return strings.Split(strings.TrimSpace(s), "\n")
}
