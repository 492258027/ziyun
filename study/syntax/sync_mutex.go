package main

import (
	"sync"
	"time"
)

//定义的类中有mutex， 那么类的函数必须传递指针（eg： d *data）, 如果传值方式，那就是两把锁。
type data struct {
	sync.Mutex
}

//注意： 类的函数一定是传指针(d *data)
func (d *data) test(s string) {
	d.Lock()
	defer d.Unlock()

	for i := 0; i < 5; i++ {
		println(s, i)
		time.Sleep(time.Second)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	var d data
	//两个匿名函数routine访问的是同一个d， 如果不是匿名函数的话， 需要传递对象d的指针。
	// func(d *data) test(s string)中一定是传指针才能保证两个routine用的是同一把锁
	go func() {
		defer wg.Done()
		d.test("write")
	}()

	go func() {
		defer wg.Done()
		d.test("read")
	}()

	wg.Wait()
}
