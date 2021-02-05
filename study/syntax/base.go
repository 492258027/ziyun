package main

import (
	"fmt"
	"reflect"
)

//变量 var开头
//常量 1预先定义的常量 true false； 2自己定义的 const开头
//变量类型： bool， 整形， string， error， 指针， 切片， 字典， 通道， 接口， 自定义的类
//切片
//字典
//切片， 字典， channel， 接口是引用语义的
func main() {
	//Var_s()
	//Const_s()
	//String_s()
	//Slice_s()
	Map_s()
}

func Var_s() {
	//变量的声明
	var a int
	//指针
	var b *int
	//string 类型
	var c string
	//error 类型
	var d error
	//切片
	var e []int
	//字典
	var f map[string]int
	//函数
	var g func(a int) int
	//自定义的类
	var h struct {
		x int
	}
	var i [10]int

	fmt.Println(reflect.TypeOf(a), reflect.TypeOf(b), reflect.TypeOf(c), reflect.TypeOf(d),
		reflect.TypeOf(e), reflect.TypeOf(f), reflect.TypeOf(g), reflect.TypeOf(h), reflect.TypeOf(i))

	//变量的赋值
	a = 10
	fmt.Println(a)
	//定义并且赋值
	j := 10
	fmt.Println(j)
}

func Const_s() {
	//true 和 false是go提前定义好的常量
	fmt.Println(reflect.TypeOf(true))

	//定义常量
	const Pi float64 = 3.14

	//枚举
	const (
		Sunday = iota
		Monday
		Tuesday
		Wednesday
	)
}

func String_s() {

	//字符串连接
	a := "hello" + "123"
	fmt.Println(a)
	//字符串的长度
	fmt.Println(len(a))
	//取字符，注意取出的字符不可以修改
	fmt.Println(a[2])
	//遍历字符串
	for i, j := range a {
		fmt.Println(i, j)
	}
}

func Slice_s() {
	//直接创建
	b := make([]int, 5, 10)
	//创建并赋值
	c := []int{1, 2, 3, 4, 5} //这一步用临时切片对象初始化c， []int是类型， 类型后面跟{}是初始化一个该类型的对象
	fmt.Println(len(b), cap(b))
	//通过数组或者切片来创建切片
	x1 := c[:]  //c的全部
	x2 := c[:2] //第二位以前的
	x3 := c[2:] //第二位以后得
	fmt.Println(x1, x2, x3)
	//动态增减
	s := []int{1, 2, 3, 4, 5}
	s = append(s, 6, 7, 8)

	e := []int{5, 6}
	s = append(s, e...)
	fmt.Println(s)
	//遍历
	for i, j := range c {
		fmt.Println(i, j)
	}
	//切片直接取值
	fmt.Println(c[0], c[1])

	//常用函数, copy 用户复制切片，如果两个切片不一样大，就按较小的容量复制
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := make([]int, 3)
	copy(slice2, slice1)
	fmt.Println(slice2)
	//字符串和切片的转换, 强制类型转换 T()
	buf := []byte("hello world!")
	fmt.Println(buf, string(buf))
}

type People struct {
	age  int
	sex  bool
	addr string
}

func Map_s() {
	//直接创建, 注意是[], 同数组
	myMap := make(map[string]People, 100)

	//创建并初始化
	myMap1 := map[string]People{
		"xuheng": People{41, true, "beijing"},
	}
	fmt.Println(myMap1)

	//元素赋值
	myMap["xuziyun"] = People{
		10, true, "suzhou",
	}

	//元素查找
	if v, ok := myMap["xuziyun"]; ok {
		fmt.Println(v)
	}

	//元素删除
	delete(myMap, "xuziyun")

	//元素的修改
	//字典不可以直接修改value, 正确的做法是返回整个value， 修改后再设置字典的键值， 或者value直接用指针类型
	//myMap1["xuziyun"].age = 12 //编译报错

	// 方法一： 返回整个value， 赋值后再设置字典的键值， 不推荐
	//m := myMap1["xuziyun"]
	//m.age =12
	//myMap1["xuziyun"] = m

	//方法二
	myMap2 := map[string]*People{
		"xuheng": &People{41, true, "beijing"},
	}
	myMap2["xuheng"].age = 42

	//map竞争访问时候需要加读写锁
}
