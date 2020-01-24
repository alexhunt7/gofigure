package master

import (
	"bytes"
	"context"
	"io/ioutil"
	"text/template"
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

func CopyFile(client *Client, src, dst, owner, group, mode string, timeoutSecs int) (*pb.FileResult, error) {
	// TODO make a proper streaming message type
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	request := &pb.FileRequest{
		Properties: &pb.FileProperties{
			Path:  dst,
			Owner: owner,
			Group: group,
			Mode:  mode,
		},
		Content: content,
	}

	return client.File(ctx, request)
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

func TemplateFile(client *Client, path, owner, group, mode, templatePath string, args interface{}, timeoutSecs int) (*pb.FileResult, error) {
	templ, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	//templName := path.Base(templatePath)
	var content bytes.Buffer
	err = templ.Execute(&content, args)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	request := &pb.FileRequest{
		Properties: &pb.FileProperties{
			Path:  path,
			Owner: owner,
			Group: group,
			Mode:  mode,
		},
		Content: content.Bytes(),
	}

	return client.File(ctx, request)
}

func LineInFile(client *Client, path, owner, group, mode, regex, line string, timeoutSecs int) (*pb.LineInFileResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	request := &pb.LineInFileRequest{
		Properties: &pb.FileProperties{
			Path:  path,
			Owner: owner,
			Group: group,
			Mode:  mode,
		},
		Regex: regex,
		Line:  line,
	}

	return client.LineInFile(ctx, request)
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
