package errCode

import (
	"errors"
)

var (
	ErrNoInstances   = errors.New("query service instance error")
	ErrNoGrpcPort    = errors.New("nofound grpc port in meta")
	ErrInvalidAToken = errors.New("access token invalid")
	ErrATokenNil     = errors.New("access token nil")
	ErrLimitExceed   = errors.New("Rate limit exceed!")
)
