package hystrix_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
	auth_pb "ziyun/auth-service/pb"
	auth_r "ziyun/auth-service/svc/client/grpc"
	"ziyun/util/bootstrap"
	consul "ziyun/util/consul"
	"ziyun/util/errCode"
)

type HystrixHandler struct {
}

func NewHystrixHandler() *HystrixHandler {
	hystrix.ConfigureCommand(bootstrap.GatewayConfig.ConsulOpStringName, hystrix.CommandConfig{
		// 设置触发最低请求阀值为 5，方便我们观察结果
		RequestVolumeThreshold: 5,
	})

	hystrix.ConfigureCommand(bootstrap.GatewayConfig.ConsulAuthName, hystrix.CommandConfig{
		// 设置触发最低请求阀值为 5，方便我们观察结果
		RequestVolumeThreshold: 5,
	})

	return &HystrixHandler{}
}

func (h *HystrixHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	reqPath := req.URL.Path
	if reqPath == "" {
		return
	}
	//按照分隔符'/'对路径进行分解，获取服务名称serviceName
	pathArray := strings.Split(reqPath, "/")

	serviceName := pathArray[1]

	// health检查
	if serviceName == "health" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	// 服务名为空, 或者服务名不是auth/opstring
	if serviceName == "" || (serviceName != bootstrap.GatewayConfig.ConsulOpStringName && serviceName != bootstrap.GatewayConfig.ConsulAuthName) {
		log.Println("serviceName invalid: ", serviceName)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	// 服务名是opstring， 需要进行鉴权， 鉴权走auth服务的grpc
	if serviceName == bootstrap.GatewayConfig.ConsulOpStringName {
		token := req.Header.Get("Atoken")
		if token == "" {
			log.Println(errCode.ErrATokenNil)
			rw.WriteHeader(http.StatusBadRequest)
			return
		} else {
			err := check_token_grpc(token)
			if err != nil {
				log.Println("check token result:", err)
				rw.WriteHeader(http.StatusForbidden)
				return
			}
		}
	}

	//参数2和参数3都是匿名函数， 参数2的作用是反向代理具体实现， 参数3的作用是接收参数2返回的错误
	err := hystrix.Do(serviceName, func() error {

		agent, err := consul.Discover(serviceName)
		if err != nil {
			return errCode.ErrNoInstances
		}
		log.Println("choice forward agent: ", agent.Address, agent.Port)

		//创建Director
		director := func(req *http.Request) {

			destPath := strings.Join(pathArray[1:], "/")

			//设置代理服务地址信息
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", agent.Address, agent.Port)
			req.URL.Path = "/" + destPath
		}

		var proxyError error

		//反向代理失败时错误处理
		errorHandler := func(ew http.ResponseWriter, er *http.Request, err error) {
			proxyError = err
		}

		proxy := &httputil.ReverseProxy{
			Director:     director,
			ErrorHandler: errorHandler,
		}

		proxy.ServeHTTP(rw, req)

		// 将执行异常反馈给hystrix
		return proxyError

	}, func(e error) error {
		//如果断路器处于open状态， 那么调用hystrix.Do时候，直接走这里
		log.Println("proxy error ", e)
		return errors.New("fallback excute")
	})

	// hystrix.Do 返回执行异常
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
	}
}

func check_token_inbody(rw http.ResponseWriter, req *http.Request) error {
	//http.request是readcloser。不能重复读取http.request里面的信息
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}

	//采用流式读取到json方式更简单，但是不会往req.Body中写
	var v interface{}
	err := json.Unmarshal(bodyBytes, &v)
	if err != nil {
		log.Println("unmarshal failure!")
		//rw.WriteHeader(http.StatusBadRequest)
		return errors.New("StatusBadRequest")
	}
	//log.Println(v)

	if m, ok := v.(map[string]interface{}); !ok { //接口的强制类型转换
		log.Println("casting to map failure!")
		//rw.WriteHeader(http.StatusBadRequest)
		return errors.New("StatusBadRequest")
	} else {
		if k, ok := m["AToken"]; !ok {
			log.Println("not found Atoke in request param!")
			//rw.WriteHeader(http.StatusBadRequest)
			return errors.New("StatusBadRequest")
		} else {
			err := check_token_grpc(k.(string)) //把interface{}转成string
			if err != nil {
				log.Println("check token result:", err)
				//rw.WriteHeader(http.StatusForbidden)
				return err
			}
		}
	}

	// 把刚刚读出来的再写进去
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}

//不需要每次新建链接， 可以复用authCli， delay修改
func check_token_grpc(token string) error {
	//鉴权服务
	agent, err := consul.Discover(bootstrap.GatewayConfig.ConsulAuthName)
	if err != nil {
		return err
	}

	var port string
	if v, ok := agent.Meta["rpcport"]; ok {
		port = v
	} else {
		return errCode.ErrNoGrpcPort
	}

	addr := agent.Address + ":" + port

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		return err
	}
	defer conn.Close()

	authCli, _ := auth_r.New(conn)
	authResp, err := authCli.Auth(context.TODO(), &auth_pb.AuthRequest{"Check_Token", "", "", "", "", token, ""})

	//鉴权未通过
	if err != nil || authResp.Valid != "true" {
		return errCode.ErrInvalidAToken
	}

	return nil
}
