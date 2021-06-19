package main

import (
	"context"
	"fmt"
	"time"
)

func doSomething(ctx context.Context) {
	select {
	case <-time.After(5 * time.Second): // 5 seconds pass
		fmt.Println("finish doing something")
	case <-ctx.Done(): // ctx is cancelled
		err := ctx.Err()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	// 创建空context的两种方法
	ctx := context.Background() // 返回一个空的context，不能被cancel，kv为空

	// todoCtx := context.TODO() // 和Background类似，当你不确定的时候使用
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(6 * time.Second)
		cancel()
	}()

	doSomething(ctx)

}
