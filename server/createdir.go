package main

import (
	pb "alex/gofigure/proto"
	"golang.org/x/net/context"
)

func (s *GofigureServer) CreateDir(ctx context.Context, cfrequest *pb.CreateDirRequest) (*pb.CreateDirReply, error) {
	return &pb.CreateDirReply{Success: true}, nil
}
