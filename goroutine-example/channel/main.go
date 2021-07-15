/**
 * @project learn-go
 * @Author 27
 * @Description //TODO
 * @Date 2021/7/15 22:46 7月
 **/
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var wg sync.WaitGroup

func player(name string, ch chan int) {
	defer wg.Done() // 任务完成
	for {
		ball, ok := <-ch // 怎样从通道拿值
		if !ok {         // 通道关闭, 说明可以理解为对手没有把球打回来
			fmt.Printf("channel is closed %s wins!\n", name)
			return
		}

		n := rand.Intn(100)
		if n%10 == 0 { // 模拟十分之一几率
			// 把球打飞，用关闭通道模拟
			close(ch)
			fmt.Printf("%s misses the ball! %s losses\n", name, name)
			return
		}
		ball++
		fmt.Printf("%s receivces ball %d\n", name, ball)
		ch <- ball
	}
}

func main() {
	wg.Add(2)

	ch := make(chan int, 0) //  unbuffered channel

	go player("heli", ch)
	go player("chong", ch)

	ch <- 0

	wg.Wait()
}
