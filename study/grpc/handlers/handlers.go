package handlers

import (
	"context"
	"errors"
	"ziyun/study/grpc/pb"
)

var (
	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

type MyString struct{}

func (s *MyString) Concat(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	return &pb.StringResponse{Result: req.A + req.B}, nil
}

func (s *MyString) Diff(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	return &pb.StringResponse{Result: req.A + req.B}, nil
}
