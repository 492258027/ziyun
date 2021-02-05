package main

import (
	"context"
	"flag"
	"fmt"
	pb "ziyun/auth-service/pb"
	r "ziyun/auth-service/svc/client/http"
)

func main() {
	flag.Parse()
	ctx := context.Background()

	authCli, _ := r.New("192.168.73.3:5150")
	AToke := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDEzODc0OTIsImlzcyI6IlN5c3RlbSIsIkNsaWVudCI6eyJDbGllbnRJZCI6IkNsaWVudElkIiwiQ2xpZW50U2VjcmV0IjoiIiwiQWNjZXNzVG9rZW5WYWxpZGl0eVNlY29uZHMiOjg2NDAwLCJSZWZyZXNoVG9rZW5WYWxpZGl0eVNlY29uZHMiOjI1OTIwMDAsIlJlZ2lzdGVyZWRSZWRpcmVjdFVyaSI6IiIsIkF1dGhvcml6ZWRHcmFudFR5cGVzIjpudWxsfSwiVXNlciI6eyJVc2VySWQiOiJVc2VySWQiLCJVc2VybmFtZSI6IlVzZXJuYW1lIiwiUGFzc3dvcmQiOiIiLCJBdXRob3JpdGllcyI6bnVsbH19.MimM5RSYzDwgDAtVPOVXZYd8zbSRkZ5dHAR3nAJ6CLQ"
	result, err := authCli.Auth(ctx, &pb.AuthRequest{"Check_Token", "", "", "", "", AToke, ""})
	if err != nil {
		fmt.Println("Check error", err.Error())
	}

	fmt.Println("result=", result)
}
