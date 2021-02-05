package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"unsafe"
)

func main() {
	song := make(map[string]interface{})
	song["name"] = "李白"
	song["timelength"] = 128
	song["author"] = "李荣浩"
	bytesData, err := json.Marshal(song)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//组request
	reader := bytes.NewReader(bytesData)
	url := "http://localhost:8080/Json"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	//发送请求
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//读body到byte[]
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	println(respBytes)

	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
}
