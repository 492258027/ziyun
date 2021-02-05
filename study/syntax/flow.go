package main

import "fmt"

//后续补充 select
func main() {

	//条件语句不需要使括号将条件包含起来。
	//在if之后，条件语句之前，可以添加变量初始化语句，使用;间隔
	if x := 100; x > 0 {
		fmt.Println(x, "x>0")
	} else if x < 0 {
		fmt.Println(x, "x<0")
	} else {
		fmt.Println(x, "x=0")
	}

	//单个case中可以出现多个结果项
	//不需要用break来明确退出一个case
	//只有在case中明确添加fallthrough关键字， 才会继续执行紧跟的下一个case, 注意fallthrough必须是case中的最后一个语句
	//可以不设定switch之后的条件表达式， 在这种情况下整个switch相当于if...else...
	i := 7
	switch i {
	case 0:
		fmt.Println("0")
	case 1:
		fmt.Println("1")
		fallthrough
	case 2:
		fmt.Println("2")
	case 3, 4, 5:
		fmt.Println("3,4,5")
	default:
		fmt.Println("default")
	}

	//for 循环遍历
	m := []int{100, 101, 102}
	for i, n := range m {
		fmt.Println(i, n)
	}

	//for 中的break和continue和goto

}
