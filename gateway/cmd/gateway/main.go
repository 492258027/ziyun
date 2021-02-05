package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"ziyun/gateway/hystrix_handler"
	"ziyun/util/bootstrap"
	consul "ziyun/util/consul"
)

//curl -X POST "http://192.168.73.3:5050/opstring/Diff" -H "Content-Type: application/json" --data '{"A":"xuheng", "B":"good"}'

//token放在body
//curl -X POST "http://192.168.73.3:5050/opstring/Diff" -H "Content-Type: application/json" --data '{"A":"xuheng", "B":"good", "AToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDYyMDMwMDEsImlzcyI6IlN5c3RlbSIsIkNsaWVudCI6eyJDbGllbnRJZCI6IkNsaWVudElkIiwiQ2xpZW50U2VjcmV0IjoiIiwiQWNjZXNzVG9rZW5WYWxpZGl0eVNlY29uZHMiOjg2NDAwLCJSZWZyZXNoVG9rZW5WYWxpZGl0eVNlY29uZHMiOjI1OTIwMDAsIlJlZ2lzdGVyZWRSZWRpcmVjdFVyaSI6IiIsIkF1dGhvcml6ZWRHcmFudFR5cGVzIjpudWxsfSwiVXNlciI6eyJVc2VySWQiOiJVc2VySWQiLCJVc2VybmFtZSI6IlVzZXJuYW1lIiwiUGFzc3dvcmQiOiIiLCJBdXRob3JpdGllcyI6bnVsbH19.HuSWATup2SVK-PiRWPw0JwDicGXy6QUilamccvplooc"}'

//token放在header
//curl -X POST "http://192.168.73.3:5050/opstring/Diff" -H "Content-Type: application/json" -H "AToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDYyMDMwMDEsImlzcyI6IlN5c3RlbSIsIkNsaWVudCI6eyJDbGllbnRJZCI6IkNsaWVudElkIiwiQ2xpZW50U2VjcmV0IjoiIiwiQWNjZXNzVG9rZW5WYWxpZGl0eVNlY29uZHMiOjg2NDAwLCJSZWZyZXNoVG9rZW5WYWxpZGl0eVNlY29uZHMiOjI1OTIwMDAsIlJlZ2lzdGVyZWRSZWRpcmVjdFVyaSI6IiIsIkF1dGhvcml6ZWRHcmFudFR5cGVzIjpudWxsfSwiVXNlciI6eyJVc2VySWQiOiJVc2VySWQiLCJVc2VybmFtZSI6IlVzZXJuYW1lIiwiUGFzc3dvcmQiOiIiLCJBdXRob3JpdGllcyI6bnVsbH19.HuSWATup2SVK-PiRWPw0JwDicGXy6QUilamccvplooc" --data '{"A":"xuheng", "B":"good"}'

//创建一个限流器
var l = rate.NewLimiter(rate.Every(1*time.Second), 2)

//中间件的方式封装
func Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//注册到consul
	consul.Register()

	//开始监听
	go func() {
		log.Println("transport", "HTTP", "addr", bootstrap.HttpConfig.Host, bootstrap.HttpConfig.Port)
		addr := bootstrap.HttpConfig.Host + ":" + strconv.Itoa(bootstrap.HttpConfig.Port)
		h := hystrix_handler.NewHystrixHandler()
		//errc <- http.ListenAndServe(addr, h)
		errc <- http.ListenAndServe(addr, Limit(h))
	}()

	//  hytrix server
	go func() {
		log.Println("transport", "hystrix", "addr", bootstrap.HystrixConfig.Host, bootstrap.HystrixConfig.Port)
		addr := bootstrap.HystrixConfig.Host + ":" + strconv.Itoa(bootstrap.HystrixConfig.Port)
		m := http.NewServeMux()
		hystrixStreamHandler := hystrix.NewStreamHandler()
		hystrixStreamHandler.Start()
		m.Handle("/hystrix/stream", hystrixStreamHandler)
		errc <- http.ListenAndServe(addr, m)
	}()

	// 开始运行，等待结束
	log.Println("exit", <-errc)
	consul.Deregister()
}
