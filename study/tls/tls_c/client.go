package main

import (
	"crypto/tls"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)

	//客户端不对服务端的证书进行CA验证, 直接通过证书解公钥，用于协商对称加密通信因子。
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hi\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	println(string(buf[:n]))
}
