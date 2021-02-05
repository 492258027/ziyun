package main

import (
	"context"
	//"database/sql/driver"
	"log"
	"sync"
	"time"
)

func main() {
	//father_ctx()
	//timeout_ctx()
	//comUse2_ctx()
	trace_ctx()
}

//ctx应用场景:
//1 超时控制: rpc调用,  restful api调用， 连接池取连接等
//2 传递请求相关的全局变量： 调用链Log， balance路由参数

//Context 提供了两个方法做初始化：
// Background 一般是所有 Context 的基础，所有 Context 的源头都应该是它。
// todo方法一般用于当传入的方法不确定是哪种类型的 Context 时，为了避免 Context 的参数为nil而初始化的 Context。

//基于 Context 派生新Context 的方法如下:
//WithCancel, WithDeadline,WithTimeout,WithValue
//这几个函数根据parent ctx派生出新的ctx，以及一个Cancel方法。 通过逐级派生， 最后会生成一颗ctx树, 如果调用了Cancel方法，对应这级Cancel的ctx开始的ctx子树中的每个ctx都会被调用对应的Cancel函数。

// WithCancel 必需要手动调用 cancel 方法，
// ctx, cancel := context.WithCancel(context.Background())

// WithDeadline 可以设置绝对时间点， 到时间点后自动调用cancel 做取消操作
// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))  //deadline的时间

// WithTimeout 是设置调用的持续时间，到指定时间后，自动调用 cancel 做取消操作。
// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)  //持续5秒

//func WithValue(parent Context, key, val interface{}) Context {}
//withValue 会构造一个新的context，新的context 会包含一对 Key-Value 数据，可以通过Context.Value(Key) 获取存在 ctx 中的 Value 值。

//Context 使用原则 和 技巧
//不要把Context放在结构体中，要以参数的方式传递，parent Context一般为Background
//应该要把Context作为第一个参数, 变量名建议都统一，如ctx。
//Context是线程安全的，可以放心的在多个goroutine中传递, sync这个包下都是线程安全的
//可以把一个 Context 对象传递给任意个数的 gorotuine，对它执行 取消 操作时，所有 goroutine 都会接收到取消信号。

var wg = sync.WaitGroup{}

//father, son, grandson 测试ctx传递性
func grandson_ctx(ctx context.Context) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		log.Println("grandson routine back!")
	}
}

func son_ctx(ctx context.Context) {
	defer wg.Done()

	wg.Add(1)
	go grandson_ctx(ctx)

	select {
	case <-ctx.Done():
		log.Println("son routine back!")
	}
}

func father_ctx() {
	ctx, cancel := context.WithCancel(context.Background())
	//ctx, _ := context.WithTimeout(context.Background(),  5*time.Second)
	//ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

	wg.Add(1)
	go son_ctx(ctx)

	time.Sleep(5 * time.Second)
	cancel() // 取消请求， 从本级ctx开始的整个子树， 都会被调用cancel

	wg.Wait()
}

//常用用法一：超时保证, 拿连接池做例子。 本级用， 超时定时器一样
func timeout_ctx() {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	p := make(chan int)

	wg.Add(2)
	go func() {
		defer wg.Done()
		// 要不就等ctx超时， 要不就等连接池返回
		select {
		case <-ctx.Done():
			// 获取链接超时了，打印错误
			println("获取链接超时：", ctx.Err())
			return
		case ret, ok := <-p:
			if !ok {
				println("连接池chan closed")
				return
			}

			// 拿到链接，校验并返回
			println("得到连接", ret)
		}
	}()

	go func() {
		defer wg.Done()
		//处理完请求的链接放回链接池
		time.Sleep(time.Second)
		p <- 1
	}()

	wg.Wait()
}

//常用用法二：判断微服务传递到这个函数之前是否已经超时
func do_sql(ctx context.Context) {
	select {
	// 校验是否已经超时，如果超时直接返回
	case <-ctx.Done():
		println("timeout")
		return
	default: //不阻塞， 没超时直接退出select
	}

	// 如果还没有超时，调用驱动做查询
	//return queryer.Query(query, dargs)
	println(" 没超时直接退出select, 执行相关访问")

}

func comUse2_ctx() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	//time.Sleep(2*time.Second)
	do_sql(ctx)
}

//链路追踪的列子
func trace_ctx() {
	// 首先从请求中拿到traceId
	// 可以把traceId 放在header里，也可以放在body中
	// 还可以自己建立一个 （如果自己是请求源头的话）
	traceId := "1234567"

	// Key 存入 ctx 中
	ctx := context.WithValue(context.Background(), "traceid", traceId)

	// 设置接口1s 超时
	ctx, _ = context.WithTimeout(ctx, time.Second)

	// query RPC 时可以携带 traceId
	req := RequestRPC(ctx)

	// 判断req， 做随后的操作
	if req == nil {
		println("ok")
	}
}

func RequestRPC(ctx context.Context) interface{} {

	// request
	select {
	// 校验是否已经超时，如果超时直接返回
	case <-ctx.Done():
		println("timeout")
		return nil
	default: //不阻塞， 没超时直接退出select
	}

	// 获取traceid，在调用rpc时记录日志
	traceId := ctx.Value("traceid") // 从ctx中取值
	println(traceId.(string))       //接口的强制转换
	return nil
}
