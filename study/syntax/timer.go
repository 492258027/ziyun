package main

import (
	"os"
	"time"
)

func main() {

	go func() {
		for {
			select {
			case <-time.After(time.Second * 3):
				println("timeout...")
				os.Exit(0)
			}
		}
	}()

	go func() {
		tick := time.Tick(time.Second) //表示每隔一段生产一条数据到tick通道中。

		for {
			select {
			case <-tick:
				println(time.Now().String())
			}
		}
	}()

	<-(chan struct{})(nil)

}

//如果明确time已经expired，并且t.C已经被取空，那么可以直接使用Reset；
//如果程序之前没有从t.C中读取过值，这时需要首先调用Stop()，
//如果返回true，说明timer还没有expir e，stop可以成功的删除timer，然后直接reset；
//如果返回false，说明stop前已经expire，需要显式drain channel。
func timerStop(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}

func timerReset(timer *time.Timer, seconds int) {
	timerStop(timer)
	timer.Reset(time.Duration(seconds) * time.Second)
}
