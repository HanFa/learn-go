/**
 * @project learn-go
 * @Author 27
 * @Description //TODO
 * @Date 2021/7/15 22:39 7月
 **/
package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func say(id string) {
	time.Sleep(time.Second)
	fmt.Println("I am done！id: " + id)
	wg. Done()  // 任务完成
}



func main() {
	wg.Add(2)  // 总共有两个任务

	go func(id string) {
		fmt.Println(id)
		wg.Done()
	}("hello")

	go say("world")

	wg.Wait()  // 等待所有任务完成 卡住，如果wg不是0
	fmt.Print("exit")
}

