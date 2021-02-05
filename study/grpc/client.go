package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"ziyun/study/grpc/pb"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:1234", grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()

	client := pb.NewStringClient(conn)
	request := &pb.StringRequest{A: "A", B: "B"}
	reply, _ := client.Concat(context.Background(), request)
	fmt.Printf("String Concat : %s concat %s = %s\n", request.A, request.B, reply.Result)
}
