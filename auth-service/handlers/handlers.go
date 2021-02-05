package handlers

import (
	"context"
	"errors"
	"ziyun/auth-service/token"

	pb "ziyun/auth-service/pb"
)

var (
	ErrNotSupportOperation = errors.New("no support operation")
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.AuthServer {
	return authService{}
}

type authService struct{}

func (s authService) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthResponse, error) {
	var resp = pb.HealthResponse{true}
	return &resp, nil
}

func (s authService) Auth(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {

	var resp pb.AuthResponse
	var atoken, rtoken, valid string
	var err error

	switch in.Type {
	case "Token":
		atoken, rtoken, err = token.Token(in.UserName, in.Passwd, in.CompanyID)
	case "Refresh_Token":
		atoken, rtoken, err = token.RefreshToken(in.RToken)
	case "Check_Token":
		valid, err = token.CheckAToken(in.AToken)
	default:
		err = ErrNotSupportOperation
	}

	resp.AToken = atoken
	resp.RToken = rtoken
	resp.Valid = valid

	return &resp, err
}
