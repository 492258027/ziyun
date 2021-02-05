package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	//json_base()
	json_stream()
	//json_unknown()
}

/*****************************************************************************************************
JSON的基本编解码
******************************************************************************************************/

type User struct {
	Name    string
	Website string
	Age     uint
	Male    bool
	Skills  []string
	Id      []byte
}

//编码函数的声明如下： func Marshal(v interface{}) ([]byte, error)
//传入参数 v 是空接口，意味着可以传入任何类型数据，如果编码成功返回对应的 JSON 格式文本，否则，通过第二个返回参数标识错误信息。
//数组和切片会转化为 JSON 里边的数组. 但[]byte 类型的值将会被转化为 Base64 编码后的字符串，slice 类型的零值会被转化为 null。
//结构体会转化为 JSON 对象，并且只有结构体里边以大写字母开头的可被导出的字段才会被转化输出，而这些可导出的字段会作为 JSON 对象的key。
//转化一个 map 类型的数据结构时，该数据的类型必须是 map[string]T（T 可以是 encoding/json 包支持的任意数据类型）。
//解码函数的声明如下：func Unmarshal(data []byte, v interface{}) error
//第一个参数是待解码的 JSON 格式文本，第二个参数表示存储解码结果的数据结构（比如上面的 User 结构体）。
func json_base() {
	user_in := User{
		"xuheng",
		"https://xueyuanjun.com",
		18,
		true,
		[]string{"Golang", "PHP", "C", "Java", "Python"},
		[]byte("1234567890"),
	}

	//编码
	u, err := json.Marshal(user_in)
	if err != nil {
		fmt.Printf("JSON 编码失败：%v\n", err)
		return
	}
	//logging.Println(u)
	log.Println(string(u)) //打印json的方式

	//解码
	var user_out User
	err = json.Unmarshal(u, &user_out)
	if err != nil {
		fmt.Printf("JSON 解码失败：%v\n", err)
		return
	}

	log.Println(user_out)

	//解码过程中， json比结构多出的字段会被丢弃， json比结构少的字段会被类型的默认值初始化
	u2 := []byte(`{"name": "xuheng", "website": "https://xueyuanjun.com", "mother": "wang"}`)
	var user2 User
	err = json.Unmarshal(u2, &user2)
	if err != nil {
		fmt.Printf("JSON 解码失败：%v\n", err)
		return
	}
	log.Println(user2)
}

/*******************************************************************************************************
JSON的流式读写
*******************************************************************************************************/
//type StringRequest struct {
//	RequestType string `json:"request_type"`
//	A           string `json:"a"`
//	B           string `json:"b"`
//}

type StringRequest struct {
	Type string
	Name string
	Pwd  string
}

//测试输入: {"Type":"diff", "Name":"xuehng", "Pwd":"123456"}
//测试输入: {"type":"diff", "name":"xuehng", "pwd":"123456"}
//json字符串中key的名字(大小写无关)必须和定义的结构体中变量名字(必须大写开头)相同才能正常解出来，
func json_stream() {
	for {
		var v StringRequest
		//读取标准输入的json， 解码json为StringRequest对象。 函数的名字都是相对于json来说的,
		if err := json.NewDecoder(os.Stdin).Decode(&v); err != nil {
			log.Println(err)
			return
		}
		log.Println(v)

		//把StringRequest对象， 编码为json格式， 然后输出到标准输出
		if err := json.NewEncoder(os.Stdout).Encode(&v); err != nil {
			log.Println(err)
		}
	}
}

/*******************************************************************************************************
解码未知结构的JSON
*******************************************************************************************************/
//允许使用 map[string]interface{} 和 []interface{} 类型的值来分别存放未知结构的 JSON 对象或数组。
//比如 JSON 对象会转换为map[string]interface{} 类型；
//测试输入： {"Name":"xuheng","Website":"https://xueyuanjun.com","Age":18,"Male":true,"Skills":["Golang","PHP","C","Java","Python"]}
func json_unknown() {

	u := []byte(`{"name": "xuheng", "website": "https://xueyuanjun.com", "age": 18, "male": true, "skills": ["Golang", "PHP"]}`)
	var v interface{} //定义一个空接口， 把json字符串解码变成map[string]interface{},  并把赋值给空接口v，
	err := json.Unmarshal(u, &v)
	if err != nil {
		fmt.Printf("JSON 解码失败：%v\n", err)
		return
	}
	fmt.Printf("JSON 解码结果: %#v\n", v)
	log.Println(v)

	if m, ok := v.(map[string]interface{}); ok { //接口的强制类型转换
		for k, s := range m {
			log.Println(k, s)
		}
	}
}
