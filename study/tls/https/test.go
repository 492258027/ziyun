package main

import (
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("https example.\n"))
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	//err := http.ListenAndServeTLS("192.168.3.115:443", "../cert/server.crt", "../cert/server.key", nil)
	err := http.ListenAndServeTLS(":443", "../cert/server.crt", "../cert/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
