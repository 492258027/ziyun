// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 7040e72f5f
// Version Date: 2020-09-19T18:42:02Z

// Package grpc provides a gRPC client for the Opstring service.
package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "ziyun/opstring-service/pb"
	"ziyun/opstring-service/svc"
)

// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn, options ...ClientOption) (pb.OpstringServer, error) {
	var cc clientConfig

	for _, f := range options {
		err := f(&cc)
		if err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}

	clientOptions := []grpctransport.ClientOption{
		grpctransport.ClientBefore(
			contextValuesToGRPCMetadata(cc.headers)),
	}
	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = grpctransport.NewClient(
			conn,
			"opstring.Opstring",
			"Health",
			EncodeGRPCHealthRequest,
			DecodeGRPCHealthResponse,
			pb.HealthResponse{},
			clientOptions...,
		).Endpoint()
	}

	var opstringEndpoint endpoint.Endpoint
	{
		opstringEndpoint = grpctransport.NewClient(
			conn,
			"opstring.Opstring",
			"Opstring",
			EncodeGRPCOpstringRequest,
			DecodeGRPCOpstringResponse,
			pb.OpstringResponse{},
			clientOptions...,
		).Endpoint()
	}

	return svc.Endpoints{
		HealthEndpoint:   healthEndpoint,
		OpstringEndpoint: opstringEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCHealthResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC health reply to a user-domain health response. Primarily useful in a client.
func DecodeGRPCHealthResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.HealthResponse)
	return reply, nil
}

// DecodeGRPCOpstringResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC opstring reply to a user-domain opstring response. Primarily useful in a client.
func DecodeGRPCOpstringResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.OpstringResponse)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCHealthRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain health request to a gRPC health request. Primarily useful in a client.
func EncodeGRPCHealthRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.HealthRequest)
	return req, nil
}

// EncodeGRPCOpstringRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain opstring request to a gRPC opstring request. Primarily useful in a client.
func EncodeGRPCOpstringRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.OpstringRequest)
	return req, nil
}

type clientConfig struct {
	headers []string
}

// ClientOption is a function that modifies the client config
type ClientOption func(*clientConfig) error

func CtxValuesToSend(keys ...string) ClientOption {
	return func(o *clientConfig) error {
		o.headers = keys
		return nil
	}
}

func contextValuesToGRPCMetadata(keys []string) grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		var pairs []string
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				pairs = append(pairs, k, v)
			}
		}

		if pairs != nil {
			*md = metadata.Join(*md, metadata.Pairs(pairs...))
		}

		return ctx
	}
}