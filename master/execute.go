package master

import (
	"context"
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

	return client.Exec(ctx, request)
}
