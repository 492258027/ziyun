package main

import (
	"fmt"
	"reflect"
)

//go中函数和变量， 小写字母开头都在本包可见， 大写字母开头的才能被其他包使用, 注意是包一级的可见性，不是文件级的可见性
func main() {
	//不定参函数的调用
	myfunc(1, 2, 3, 4)

	//任意类型不定参函数
	var v1 int = 1
	var v2 string = "xuheng"
	myPrintf(v1, v2)

	//测试返回值
	fmt.Println(myfunc3(1, 2))
	fmt.Println(myfunc4(1, 2))

	//测试匿名函数
	myfunc5()

	//闭包是指函数及其引用的变量的组合体， 我们常用在匿名函数上, 使匿名函数可以访问匿名函数外的变量。
	//我们常用的defer函数就是一个闭包， defer内部可以访问闭包外的变量
	myfunc6()
	myfunc7()

	myfunc8()
}

//不定参函数及其传递
func myfunc2(args ...int) {
	fmt.Println(reflect.TypeOf(args))
}

func myfunc(args ...int) {
	//相当于把1,2,3这样的可变参组成切片， 函数内部直接当切片用， ...就是切片的语法糖
	for _, arg := range args {
		fmt.Println(arg)
	}

	//如果继续往下传递，需要切片散列开
	myfunc2(args...)
}

//任意类型的不定参函数
func myPrintf(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case int:
			fmt.Println(arg, "is a int value")
		case string:
			fmt.Println(arg, "is a string value")
		default:
			fmt.Println(arg, "is a unknow type")
		}
	}
}

//匿名返回值:只有数据类型如, 如(int,int)    此时,函数体中要显示返回return  1,1  或return  a,b
func myfunc3(a, b int) (int, error) {
	return a + b, nil
}

//命名返回值:既有返回值类型,又有返回值变量的名称.如:(res int)    此时,函数体中不用显示返回,如return就可以.
func myfunc4(a, b int) (req int, err error) {
	req = a + b
	//err = nil
	return
}

//匿名函数作为参数
func test1(f func()) {
	f()
}

//匿名函数作为返回值
func test2() func(int, int) int {
	return func(x, y int) int {
		return x + y
	}
}

func myfunc5() {

	//方式一: 匿名函数可以赋值给一个变量
	f := func(x, y int) int {
		return x + y
	}
	fmt.Println(f(4, 5))

	//方式二: 匿名函数可以直接执行
	func(x int) {
		fmt.Println(x)
	}(4)

	//方式三： 匿名函数做为参数
	test1(func() {
		fmt.Println("hello world!")
	})

	//方式四: 匿名函数做为返回值
	k := test2()
	fmt.Println(k(1, 2))
}

//测试闭包
func test4(x int) func() {
	return func() {
		fmt.Println(x)
	}
}

func myfunc6() {
	f := test4(123)
	f()
}

func myfunc7() {
	j := 5
	f := func() func() {
		i := 10
		return func() {
			fmt.Println(i, j)
		}
	}()

	f()
	//defer f()

	j = 10

	f()
	//defer f()
	//闭包可以访问闭包外的变量j, 取的是闭包执行时变量j的值。
	//i值被隔离， 只能在闭包内访问，闭包外不能访问。
	//加上defer以后，因为延迟调用， 访问到的j是最后赋值时的值10
}

//函数前面加上defer延迟执行。
//如果对状态敏感， 常使用闭包或者指针
func myfunc8() {
	x, y := 1, 2

	//defer后是闭包， 其中a是匿名函数传参， 取的是延迟函数压栈时的x的当前值， y是闭包中引用匿名函数外的变量， 会取延迟函数函数调用时y的当前值
	defer func(a int) {
		println("defer x, y =", a, y)
	}(x)

	x += 100
	y += 100
	print(x, y)
}
