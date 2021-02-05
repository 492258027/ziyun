package main

import (
	"sync"
	"time"
)

//sync.WaitGroup
func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1) //最外层等结果的routine来设置累加计数,  add一次起一个routine

		go func(id int) {
			defer wg.Done() //每个routine内部递减计数

			time.Sleep(time.Second)
			println("goroutine", id, "done")
		}(i)
	}

	println("main......")
	wg.Wait() //最外层等结果的routine来调用wait(), 调用后阻塞， 直到计数器归零
	println("main exit")
}
