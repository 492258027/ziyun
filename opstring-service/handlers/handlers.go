package handlers

import (
	"context"
	pb "ziyun/opstring-service/pb"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.OpstringServer {
	return opstringService{}
}

type opstringService struct{}

func (s opstringService) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthResponse, error) {
	var resp = pb.HealthResponse{true}
	return &resp, nil
}

func (s opstringService) Opstring(ctx context.Context, in *pb.OpstringRequest) (*pb.OpstringResponse, error) {
	var req string
	switch in.Type {
	case "Diff":
		req = in.A + "-" + in.B
	case "Concat":
		req = in.A + "+" + in.B
	default:
		req = "invalid request parameter"
	}

	var resp = pb.OpstringResponse{req}
	return &resp, nil
}
