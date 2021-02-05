package logging

import (
	"log"
	"time"
	"ziyun/study/middleware/service"
)

//类封装的方式增加log
type Logmw struct {
	Logger *log.Logger
	Next   service.IStringServicer
}

func (mw Logmw) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Println(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Next.Uppercase(s)
	return
}

func (mw Logmw) Count(s string) (n int) {
	defer func(begin time.Time) {
		mw.Logger.Println(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.Next.Count(s)
	return
}

//对外提供优雅的方式初始化logmw
func LoggingMiddlewareElegant(logger *log.Logger) service.ServiceMiddleware {
	return func(next service.IStringServicer) service.IStringServicer {
		return Logmw{logger, next}
	}
}
