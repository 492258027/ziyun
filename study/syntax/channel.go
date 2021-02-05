package main

import (
	"fmt"
	"sync"
	"time"
)

//创建, 容量
//通道值就是指向地址的指针或者nil， 当为地址时候， 直接赋值，引用传递。 对nil的chan的读写只有在select编译内部不报错, 因为无论读写都是阻塞的，所以select也不会处理这段。
//单向通道
//通道安全的收发模式， close对通道读写的影响， close一般由send的routine调用， 通道的选择 select， 及default的注意事项
//工厂模式

func main() {
	//singleChan()
	//okRecv()
	//rangeRecv()
	//closeChan()
	//selectChan1()
	//selectChan2()
	//factory()
	//nilChan()
	writeblock()
}

// 通道的创建，容量， nil
func baseChan() {
	//c1:= make(chan struct{})
	c1 := make(chan int)
	//对于阻塞的channel， len和cap的返回值都是0, 可以通过这种方式判断是否是同步通道
	fmt.Println(len(c1), cap(c1))

	c2 := make(chan int, 10)
	//对于非阻塞的channel， len返回实际占用的长度，cap返回容量
	fmt.Println(len(c2), cap(c2))

	//通道值是一个具体的地址， 0x值或者是nil
	fmt.Println(c1)

	c2 = nil
	fmt.Println(c2)
	//c2 <- 1 //写入会报fatal， 具体怎么安全的写现在不知道 delay
}

//单向通道： 通常使用类型转换来获取单向通道
//关闭单向通道时候， 要不就关闭类型转换之前的原始通道， 要不就关闭写方向通道。
func singleChan() {
	var wg sync.WaitGroup
	wg.Add(2) //一般都在go外面主进程add

	c := make(chan int)
	var send chan<- int = c //通道引用传递的
	var recv <-chan int = c

	go func() {
		defer wg.Done()
		for {
			x, ok := <-recv
			if !ok {
				return
			}
			println(x)
		}
	}()

	go func() {
		defer wg.Done()
		//defer close(c)   //注意是发送端close通道
		defer close(send) //发送端关闭写方向的通道是可以的，关闭recv会报错

		for i := 0; i < 3; i++ {
			send <- i
		}
	}()

	wg.Wait()
}

//安全收发方式
//发就是直接发， 无保护， 因为发到关闭状态的chan直接panic， 所以close命令由发方来调用
//接有两种方式， 1 ok-idom  3 range
func okRecv() {
	var wg sync.WaitGroup
	wg.Add(1)

	c := make(chan int)

	go func() {
		defer wg.Done()
		for {
			x, ok := <-c
			if !ok { //通道关闭
				println("channel closed!")
				return
			}
			println(x)
		}
	}()

	c <- 1
	c <- 1

	close(c)

	wg.Wait()
}

func rangeRecv() {
	var wg sync.WaitGroup
	wg.Add(1)

	c := make(chan int)

	go func() {
		defer wg.Done()
		for x := range c { //一直阻塞等待接收消息， 直到缓存为空并且通道关闭
			println(x)
		}
	}()

	c <- 1
	time.Sleep(time.Second * 5)
	c <- 1

	close(c)

	wg.Wait()
}

//向已关闭的通道发送数据， 引发panic, 因为写会引起panic， 所有一般都是写的routine来调用close函数关闭channel
//从已关闭的通道接收数据，返回已缓冲数据或者零值
//无论收发， nil通道都会阻塞
func closeChan() {
	c := make(chan int, 3)
	c <- 10
	c <- 20

	close(c)

	//可以正常读出来， 结束条件：缓存为空并且通道关闭
	for x := range c {
		println(x)
	}
	//同上， 效果一样的
	//for {
	//	x, ok:= <-c
	//	if !ok{
	//		return
	//	}
	//	println(x)
	//}
}

//如果要同时处理多个通道，可以选用select，他会随机选择一个可用的通道做收发操作(单独通道也可以用select)
//select中可以加default， 这样select就不阻塞了， 死循环一直跑
func selectChan1() {
	var wg sync.WaitGroup
	wg.Add(2)

	a := make(chan int)
	b := make(chan int)

	go func() {
		defer wg.Done()
		for {
			select {
			case x, ok := <-a:
				if !ok {
					println("channel:", a, "closed")
					//break  //死循环， 跳出了select
					return
				}
				println(x)
			case x, ok := <-b:
				if !ok {
					println("channel:", b, "closed")
					//break //死循环， 跳出了select
					return
				}
				println(x)
			}
		}
	}()

	//go func() {
	//	defer wg.Done()
	//	defer close(a)
	//	defer close(b)
	//	for i := 0; i < 10; i++ {
	//		select {   //随机选一个通道做发送
	//			case a <- i:
	//			case b <- i:
	//		}
	//	}
	//}()
	//
	//wg.Wait()

	//崩溃, wait检测到没有读的routine了， 不用wg就没事， 但是还是推荐使用wg
	go func() {
		defer wg.Done()
		defer close(a)

		for i := 0; i < 10; i++ {
			a <- i
		}
	}()

	go func() {
		defer wg.Done()
		defer close(b)

		for i := 10; i < 15; i++ {
			b <- i
		}
	}()

	wg.Wait()
}

//对nil的chan的读写只有在select编译内部不报错, 因为无论读写都是阻塞的，所以select也不会处理这段。
//另外注意break是跳出select
func selectChan2() {
	var wg sync.WaitGroup
	wg.Add(3)

	a := make(chan int)
	b := make(chan int)

	go func() {
		defer wg.Done()
		for {
			select {
			case x, ok := <-a:
				if !ok {
					println("channel:", a, "closed")
					a = nil
					break
				}
				println(x)
			case x, ok := <-b:
				if !ok {
					println("channel:", b, "closed")
					b = nil
					break //不是用return， 只是跳出了select
				}
				println(x)
			}
			if a == nil && b == nil {
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		defer close(a)

		for i := 0; i < 10; i++ {
			a <- i
		}
	}()

	go func() {
		defer wg.Done()
		defer close(b)

		for i := 10; i < 15; i++ {
			b <- i
		}
	}()

	wg.Wait()
}

//工厂模式
type receiver struct {
	sync.WaitGroup
	data chan int
}

func NewReceiver() *receiver {
	r := &receiver{
		data: make(chan int),
	}

	r.Add(1)

	go func() {
		defer r.Done()
		for x := range r.data {
			println("recv:", x)
		}
	}()

	return r
}

func factory() {
	r := NewReceiver()
	r.data <- 1
	r.data <- 2

	close(r.data) //写入方调用close

	r.Wait()
}

//插入nil，也占一个容量
func nilChan() {
	c := make(chan *int, 10)
	fmt.Println(len(c), cap(c))
	a := 5
	c <- &a
	c <- nil
	fmt.Println(len(c), cap(c))
}

//写入阻塞测试
func writeblock() {
	//判断微服务传递到这个函数之前是否已经超时
	c := make(chan int, 1)
	c <- 1

	select {
	case c <- 2:
		println("write ok!")
	case <-time.After(time.Second * 3):
		println("write timeout")
	}
}
