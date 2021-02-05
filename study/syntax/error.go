package main

import (
	"errors"
	"log"
)

// 标准库把error定义为接口类型
// type error interface{
//		Error() string
// }

//标准库提供方法，可以方便的创建简单错误文文的error对象， 这个对象因为实现了Error(), 可以赋值给接口error
func div(x, y int) (int, error) {
	if y == 0 {
		return 0, errors.New("division by zero") //对象赋值给接口
	}

	return x / y, nil //接口的默认值是nil
}

func main() {

	//panic会立即中断当前函数的流程，执行defer调用
	//defer调用中recover可以捕获并返回panic提交的错误对象
	//recover必须在defer函数中才可以正常工作
	//中断的错误会沿调用栈向外层传递， 要不被外层捕获， 要不导致进程崩溃

	defer func() {
		if err := recover(); err != nil {
			//debug.PrintStack() //调用栈调试
			log.Fatalln("recover：", err)
		}
	}()

	z, err := div(5, 0)
	if err != nil {
		//logging.Fatalln(err) //注意打印err中文字的方法
		panic(err)
	}

	println(z)
}
