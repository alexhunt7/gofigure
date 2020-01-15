package master

import (
	"context"
	"fmt"
	"time"

	pb "github.com/alexhunt7/gofigure/proto"
)

func Exec(client *Client, args []string, timeoutSecs int) (*pb.ExecResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	request := &pb.ExecRequest{
		Executable: args[0],
		Args:       args[1:],
	}

	res, err := client.Exec(ctx, request)
	if err == nil && int(res.ReturnCode) != 0 {
		return res, fmt.Errorf("non-zero return code: %d\n\nSTDOUT:\n%s\n\nSTDERR:\n%s", res.ReturnCode, res.Stdout, res.Stderr)
	}
	return res, err
}
