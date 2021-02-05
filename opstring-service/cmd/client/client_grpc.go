package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"time"
	pb "ziyun/opstring-service/pb"
	r "ziyun/opstring-service/svc/client/grpc"
)

func main() {
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:5040", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}
	defer conn.Close()

	strCli, _ := r.New(conn)
	result, err := strCli.Health(ctx, &pb.HealthRequest{})
	if err != nil {
		fmt.Println("Check error", err.Error())
	}

	fmt.Println("result=", result)
}
