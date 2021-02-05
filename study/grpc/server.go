package main

import (
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	r "ziyun/study/grpc/handlers"
	"ziyun/study/grpc/pb"
)

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//初始化对象
	str := new(r.MyString)

	// Creates a new gRPC server
	grpcServer := grpc.NewServer()

	//把server和对象绑定
	pb.RegisterStringServer(grpcServer, str)
	//pb.RegisterStringServer(grpcServer, &r.MyString{})

	//开始工作
	grpcServer.Serve(l)
}
