package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

var addr = flag.String("addr", "192.168.73.3:5360", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	h := mux.NewRouter()
	h.HandleFunc("/Echo", Echo)
	//http.ListenAndServe("192.168.73.3:6060", h)
	//http.ListenAndServe(*addr, h)
	http.ListenAndServe(":6060", h)
}

func Echo(w http.ResponseWriter, r *http.Request) {

	params, _ := url.ParseQuery(r.URL.RawQuery)
	if v, ok := params["name"]; ok {
		log.Println(v[0])
	}
	if v, ok := params["pwd"]; ok {
		log.Println(v[0])
	}

	//升级到websocket处理方式
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		http.NotFound(w, r)
		return
	}

	defer conn.Close()

	//设置读取消息大小上线
	conn.SetReadLimit(1024)

	//解析参数
	//此处可以解析http相关参数

	//消息收发处理
	for {
		//从client端接收消息
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Println("recv: ", string(message))

		//给client回Ack
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
