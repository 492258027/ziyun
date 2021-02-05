package main

import (
	"fmt"
	"math/rand"
)

//对象赋值接口: 对象只要实现了接口声明的函数即可赋值
//接口赋值接口: 接口如果包含其他接口声明的函数即可把自己赋值给其他接口， 往子集赋， 典型的例子任何接口都可以给空接口赋值
//接口指向对象的类型查询
//接口组合， 接口为引用语义， 接口的匿名组合的时候不用加*。 接口组合不能有同名的方法。
//接口的默认值 nil

type AccessMysql struct {
	endpoint []string
}

//随机访问
func (p *AccessMysql) balance() {
	i := rand.Intn(len(p.endpoint))
	fmt.Println(i, p.endpoint[i])
}

type AccessMq struct {
	endpoint []string
	pos      int
}

//轮询访问
func (p *AccessMq) balance() {
	i := p.pos / len(p.endpoint)
	fmt.Println(i, p.endpoint[i])
}

//接口通常以er结尾
type IAccesser interface {
	balance()
}

func route(a []IAccesser) {
	for _, v := range a {
		if m, ok := v.(*AccessMysql); ok { //把接口指向的对象强制转换成指定类型，然后判断返回值， 确定是否可以强制转换
			fmt.Println("type AccessMysql")
			m.balance()
		} else if m, ok := v.(*AccessMq); ok {
			fmt.Println("type AccessMq")
			m.balance()
		} else {
			fmt.Println("type unknown")
		}
	}
}

func route2(a []IAccesser) {
	for _, v := range a {
		switch v.(type) { //接口指向对象的类型
		case *AccessMysql:
			println("type AccessMysql")
			v.balance()
		case *AccessMq:
			println("type AccessMq")
			v.balance()
		default:
			println("type unknown")
		}
	}
}

func main() {
	mysql := AccessMysql{[]string{"192.168.1.1", "192.168.1.2"}}
	mq := AccessMq{[]string{"127.0.0.1", "127.0.0.2"}, 0}
	//route([]IAccesser{ &mysql, &mq})  //注意接口赋值给的对象的地址
	route2([]IAccesser{&mysql, &mq}) //注意接口赋值给的对象的地址
}
