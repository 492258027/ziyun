package main

import (
	"fmt"
	"log"
	"os"
)

//切片， 字典， channel， 接口是引用语义的

//定义类
type Persion struct {
	age  int
	name string
}

//约定俗成的初始化方法，  初始化对象时候显示调用
func NewPersion(age int, name string) *Persion {
	return &Persion{age, name}
}

//直接写p也可以，写p*也可以， 我们一般写指针的方式。
//go给指针赋值是需要取地址&, 使用指针时直接打点就行
func (p Persion) getAge() int {
	return p.age
}

func (p *Persion) getName() string {
	return p.name
}

//测试匿名组合
type Persion_s struct {
	age         int
	name        string
	*log.Logger //匿名组合
}

//约定俗成的初始化方法
func NewPersion_s(age int, name string, l *log.Logger) *Persion_s {
	return &Persion_s{age, name, l}
}

func main() {
	//Persion对象
	p := NewPersion(10, "xuziyun")       //返回的是对象的指针
	fmt.Println(p.getAge(), p.getName()) //指针直接打点就行

	//Persion_s对象
	p_s := NewPersion_s(10, "xuziyun", log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime)) //返回的是对象的指针
	p_s.Println(p.name, p.age)

	//测试显示字段名
	mytest1()
}

//需要注意的是go实现最小的面向对象， 匿名嵌入不是继承， 无法实现多态。 它更倾向于组合优于继承这种思想， 它将模块分解成相互独立的更小单元,
//分别处理不同方面的的需求，最后以匿名嵌入的方式组合到一起， 共同实现对外的接口。组合没有父子依赖， 整体和局部松耦合，可任意增加来实现扩展，
//各个单元职责单一， 互无关联， 实现和维护更加简单。 尽管接口也是多态的一种实现形式， 但是我认为应该和基于继承体系的多态分离开来。
//编译器从最外的显示命名字段开始， 逐步向内查找匿名字段，如匿名字段被外层同名字段屏蔽， 则必须使用显示字段名
//如果多个相同层级的匿名字段成员重名， 也只能使用显示字段名访问
type file struct {
	name string
}

type data struct {
	name string
	file
}

func mytest1() {
	d := data{
		"data",
		file{"file"},
	}

	log.Println(d.name, d.file.name)
}
