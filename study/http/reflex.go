package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {

	//创建Director
	director := func(req *http.Request) {
		//查询原始请求路径
		reqPath := req.URL.Path
		if reqPath == "" {
			return
		}
		//按照分隔符'/'对路径进行分解，获取服务名称serviceName
		pathArray := strings.Split(reqPath, "/")
		serviceName := pathArray[1]
		println(serviceName)

		//重新组织请求路径，去掉服务名称部分
		destPath := strings.Join(pathArray[2:], "/")

		ServiceAddress := "localhost"
		ServicePort := 8080

		//设置代理服务地址信息
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", ServiceAddress, ServicePort)
		req.URL.Path = "/" + destPath
	}

	proxy := &httputil.ReverseProxy{Director: director}

	//开始监听
	http.ListenAndServe(":9090", proxy)
}
