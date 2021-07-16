/**
 * @project learn-go
 * @Author 27
 * @Date 2021/7/15 23:55 7月
 **/
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var counter int32
var mtx sync.Mutex
var ch = make(chan int, 1)
/*
这里改为 bufferd 原因是最后一个循环把数送进ch后没有再一个循环从ch拿值
加上bufferd 最后ch存有一个值， 这样可以结束goroutine 外面从里面拿出那个结果
 */

// ================================== 竞争条件情况
func UnsafeIncCounter() {
	defer wg.Done()
	for i := 0; i < 10000; i ++ {
		//counter ++
		// 以上实际上是以下逻辑步骤
		temp := counter  // 读
		temp = temp + 1  // 计算
		counter = temp  // 写入
	}
}

func raceConditionDemo() {
	wg.Add(2)

	go UnsafeIncCounter()
	go UnsafeIncCounter()

	wg.Wait()

	fmt.Println(counter)
}

// ================================== 正常且安全情况
func MutexIncCounter() {
	defer wg.Done()
	for i := 0; i < 10000; i ++ {
		// race condition here 竞争条件
		mtx.Lock()
		counter ++  // 共享资源
		mtx.Unlock()
	}
}

func AtomicIncCounter() {
	defer wg.Done()
	for i := 0; i < 10000; i ++ {
		// race condition here 竞争条件
		atomic.AddInt32(&counter, 1)
	}
}

func mutexDemo() {
	wg.Add(2)

	go MutexIncCounter()
	go MutexIncCounter()

	wg.Wait()

	fmt.Println(counter)
}

func atomicDemo() {
	wg.Add(2)

	go AtomicIncCounter()
	go AtomicIncCounter()

	wg.Wait()

	fmt.Println(counter)
}

func ChannelIncCounter() {
	defer wg.Done()
	for i := 0; i < 10000; i ++ {
		count := <- ch
		count ++
		fmt.Println(count)
		ch <- count
	}
}

func channelDemo() {
	wg.Add(2)

	// race condition
	go ChannelIncCounter()
	go ChannelIncCounter()

	ch <- 0

	wg.Wait()

	fmt.Println(<- ch)
}


func main() {
	//  运行以下是race condition 情况
	//raceConditionDemo()

	// 以下是正常情况
	//mutexDemo()

	//atomicDemo()

	channelDemo()
}
