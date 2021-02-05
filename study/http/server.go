package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

//curl -X POST 'http://localhost:10098/oauth/token?grant_type=password' \
//-H 'Authorization: Basic Y2xpZW50SWQ6Y2xpZW50U2VjcmV0' \
//-H 'Content-Type: multipart/form-data' \
//-F username=simple \
//-F password=123456

/////////////////////////////////////////content-type/////////////////////////////////////////////
/* query形式
当发起一次POST请求时，若未指定content-type，则默认content-type为application/x-www-form-urlencoded。
即参数会以Form Data的形式进行传递，不会显式出现在请求url中。
curl -X POST "http://10.17.8.114:8080/form?type=Diff"

curl使用-d参数以后，HTTP 请求会自动加上标头Content-Type : application/x-www-form-urlencoded。
curl -X POST "http://10.17.8.114:8080/form" -d "type=Diff"

以上两种都可以用form来处理
参见FormHandler函数例子
*/

/* json
curl -X POST "http://10.17.8.114:8080/Json?type=Diff" -H "Content-Type: application/json" --data '{"userID":"10001"}'
curl -X POST "http://10.17.8.114:8080/Json?type=Diff" -H "Content-Type: application/json" --data '{"Name":"xuheng","Website":"https://xueyuanjun.com","Age":18,"Male":true,"Skills":["Golang","PHP","C","Java","Python"]}'
参见JsonHandler函数例子
*/

/* multi form
curl 的-F 发送文件
curl -X POST  'http://10.17.8.114:8080/multipart' -H 'Content-Type: multipart/form-data' -F 'payload={"aa":"bb"}' -F 'file1=@a1.txt' -F 'file2=@a2.zip'
参见MultataHandler函数例子
*/

/*
OperatorHandler 示例
curl -X POST "http://192.168.73.3:8080/operator/Diff" -H "Content-Type: application/json" -d '{"Username":"xuheng","Password":"123456"}'
*/

/////////////////////////////////////////mux路由/////////////////////////////////////////////

/* mux
普通路由， 参数路由， Matching路由， 分组路由， 路由中间件
*/

func main() {
	r := mux.NewRouter()
	//以下测试Content-Type
	r.HandleFunc("/form", FormHandler).Methods("POST")
	r.HandleFunc("/Json", JsonHandler).Methods("POST")
	r.HandleFunc("/multipart", MultiHandler).Methods("POST")
	//以下测试mux的路由
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/hello/haha", HelloHandler)
	r.HandleFunc("/products/{key:[a-z,A-Z]+}", ProductsHandler).Host("10.17.8.114").Schemes("http").Methods("POST")
	//以下测试swagger
	r.HandleFunc("/operator/{key:[a-z,A-Z]+}", OperatorHandler)

	//设置端口 路由
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Handler:      r,
	}
	//启动监听
	server.ListenAndServe()
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "index, index")
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello handle!")
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if v, ok := vars["key"]; ok {
		println(v)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello, Products")
}

//必须大写开头，必须和json参数的key(不区分大小写)一样才能解
type user struct {
	Username string
	Password string
}

type out struct {
	Message string
	Code    int
}

func OperatorHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var u user
	var o out

	if requestType, ok := vars["key"]; ok {
		//把输入json串转换成user结构体
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			log.Println(err)
			return
		}

		switch requestType {
		case "Concat":
			o.Message = u.Username + "+" + u.Password
			o.Code = 200
		case "Diff":
			o.Message = u.Username + "-" + u.Password
			o.Code = 200
		default:
			o.Message = "invalid request parameter"
			o.Code = 500
		}

		//返回json串给客户端
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(o); err != nil {
			log.Println(err)
			return
		}
	}
}

func FormHandler(w http.ResponseWriter, r *http.Request) {

	//打印header参数
	fmt.Println("-----------------header begin-----------------")
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			fmt.Printf("%s=%s\n", k, v[0])
		}
	}

	fmt.Println("-----------------form begin-----------------")
	r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form { //遍历map
			fmt.Printf("%s=%s\n", k, v[0])
		}
	}

	//fmt.Println("-----------------query begin-----------------")
	//params,_ := url.ParseQuery(r.URL.RawQuery)
	//msg := params["type"][0]
	//fmt.Println(msg)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello, Products")
}

func JsonHandler(w http.ResponseWriter, r *http.Request) {

	//打印header参数
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			fmt.Printf("%s=%s\n", k, v[0])
		}
	}
	fmt.Println("-----------------header end-----------------")

	// 解客户端发来的json， 因为是未知结构的JSON， 用map[string]interface{}来接
	m := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&m)
	for key, value := range m { // 遍历map
		log.Println("key:", key, "value :", value)
	}

	// 返回json字符串给客户端
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func MultiHandler(w http.ResponseWriter, r *http.Request) {

	//打印header参数
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			fmt.Printf("%s=%s\n", k, v[0])
		}
	}
	fmt.Println("-----------------header end-----------------")

	err := r.ParseMultipartForm(1048576)
	if err != nil {
		log.Printf("Cannot ParseMultipartForm, error: %v\n", err)
		return
	}

	if r.MultipartForm == nil {
		log.Printf("MultipartForm is null\n")
		return
	}

	if r.MultipartForm.Value != nil {
		parseMultipartFormValue(r.MultipartForm.Value)
	}

	if r.MultipartForm.File != nil {
		parseMultipartFormFile(r, r.MultipartForm.File)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello, Products")
}

// parse form data
func parseMultipartFormValue(formValues map[string][]string) {
	for formName, values := range formValues {
		log.Printf("Value formname: %s\n", formName)
		for i, value := range values {
			log.Printf("      formdata[%d]: content=[%s]\n", i, value)
		}
	}
	return
}

// parse form file
func parseMultipartFormFile(r *http.Request, formFiles map[string][]*multipart.FileHeader) {
	for formName, _ := range formFiles {
		// func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
		// FormFile returns the first file for the provided form key
		formFile, formFileHeader, _ := r.FormFile(formName)

		log.Printf("File formname: %s, filename: %s, file length: %d\n", formName, formFileHeader.Filename, formFileHeader.Size)

		if strings.HasSuffix(formFileHeader.Filename, ".zip") {
			zipReader, _ := zip.NewReader(formFile, formFileHeader.Size)
			for i, zipMember := range zipReader.File {
				f, _ := zipMember.Open()
				defer f.Close()

				buf, _ := ioutil.ReadAll(f)
				log.Printf("     formfile[%d]: filename=[%s], size=%d, content=[%s]\n", i, zipMember.Name, len(buf), strings.TrimSuffix(string(buf), "\n"))
			}
		} else {
			var b bytes.Buffer
			_, _ = io.Copy(&b, formFile)
			log.Printf("     formfile: content=[%s]\n", strings.TrimSuffix(b.String(), "\n"))
		}
	}
}
