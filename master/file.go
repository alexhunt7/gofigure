package master

import (
	"context"
	"time"

	pb "github.com/alexhunt7/gofigure/proto"
)

func Stat(client *Client, path string, timeoutSecs int) (*pb.StatResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	request := &pb.FilePath{
		Path: path,
	}

	return client.Stat(ctx, request)
}

func File(client *Client, path, owner, group, mode string, content []byte, timeoutSecs int) (*pb.FileResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	request := &pb.FileRequest{
		Properties: &pb.FileProperties{
			Path:  path,
			Owner: owner,
			Group: group,
			Mode:  mode,
		},
		Content: content,
	}

	return client.File(ctx, request)
}

func Directory(client *Client, path, owner, group, mode string, timeoutSecs int) (*pb.DirectoryResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	request := &pb.FileRequest{
		Properties: &pb.FileProperties{
			Path:  path,
			Owner: owner,
			Group: group,
			Mode:  mode,
		},
	}

	return client.Directory(ctx, request)
}

func Delete(client *Client, path string, recursive bool, timeoutSecs int) (*pb.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	request := &pb.DeleteRequest{
		Path:      path,
		Recursive: recursive,
	}

	return client.Delete(ctx, request)
}
