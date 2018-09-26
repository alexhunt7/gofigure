package main

import (
	pb "alex/gofigure/proto"
	"golang.org/x/net/context"
)

func (s *GofigureServer) CreateFile(ctx context.Context, cfrequest *pb.CreateFileRequest) (*pb.CreateFileReply, error) {
	return &pb.CreateFileReply{Success: true}, nil
}
