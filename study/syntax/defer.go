package main

import "log"
import "errors"

//return是分为两步执行的，第一步赋值给返回值(return 可以是表达式)，第二步真正的返回到函数外部。而defer是在第一步之后执行。
//“return i”把 i赋值给返回值i（当然，这里return的值就是i，所以其实没有赋值），此时i=1,然后再执行 defer，i=2,返回的i最终值是2。
func func1() (i int) {
	i = 1

	defer func() {
		i++
	}()

	return i
}

//“return i”把i=1赋值给返回值，但是这里的返回值没有显示声明，会生成一个临时变量，假设叫‘tmp’，即tmp=1。然后，执行defer的时候，i=2。但是这个和'tmp'没关系。所以最终返回的是1。
func func2() int {
	i := 1

	defer func() {
		i++
	}()

	return i
}

//“return &i”把i的地址赋值给返回值tmp。然后执行defer的时候，对tmp指向的那个值修改成了2。所以结果是2。
func func3() *int {
	i := 1
	defer func() {
		i++
	}()

	return &i
}

//综上，如果是显式命名的返回值，则defer中可以对其操作。
//如果是非显式命名的返回值，则返回时会新定义一个返回变量，defer操作不到。

//此函数区分panic和fatal的区别
func func4() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recover:", err)
		}
	}()

	//函数停止向下执行，执行defer，一层一层向上抛panic，直到recover捕获到。
	panic("mod new fatal")

	//打印输出内容,应用程序马上退出(os.Exit)，defer函数不会执行
	log.Fatalln("mod new fatal")
}

//此函数尝试在recover捕获错误的情况下，也能把err带出来, 必须显式声明err才行
func func5() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recover:", err)
		}
	}()

	err = errors.New("i found err")

	panic("new fatal")

	return err
}

func main() {
	rq := func5()
	if rq != nil {
		log.Println(rq.Error())
	} else {
		log.Println("return nil")
	}

	//	log.Println(func1())
	//	log.Println(func2())
	//	log.Println(*func3())
}
