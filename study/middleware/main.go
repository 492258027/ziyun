package main

import (
	"log"
	"os"
	"ziyun/study/middleware/logging"
	"ziyun/study/middleware/service"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	var svc service.IStringServicer

	//对封装一层的logging， 直接初始化并赋值给接口svc
	//svc = logging.Logmw{logger, service.StringService{}}
	//svc.Uppercase("xuheng")

	//对封装一层的logging，提供优雅的方式来初始化logmw， 并赋值给接口svc
	svc = service.StringService{}
	svc = logging.LoggingMiddlewareElegant(logger)(svc)
	svc.Uppercase("xuheng") //这里需要注意下： 只能访问到最外层函数。
}
