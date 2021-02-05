package main

import (
	"context"
	"flag"
	"fmt"
	pb "ziyun/opstring-service/pb"
	r "ziyun/opstring-service/svc/client/http"
)

func main() {
	flag.Parse()
	ctx := context.Background()

	strCli, _ := r.New("192.168.73.3:5250")

	result, err := strCli.Opstring(ctx, &pb.OpstringRequest{"Diff", "xuheng", "good"})
	if err != nil {
		fmt.Println("Check error", err.Error())
	}

	fmt.Println("result=", result)
}
