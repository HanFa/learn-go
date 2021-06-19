package main

import (
	"sync"
	"sync/atomic"
)

func atomicIncCounter(counter *int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		atomic.AddInt64(counter, 1)
	}
}

func mutexIncCounter(counter *int64, wg *sync.WaitGroup, mtx *sync.Mutex) {
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		mtx.Lock()
		*counter++
		mtx.Unlock()
	}
}

func ConcurrentAtomicAdd() int64 {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var counter int64 = 0
	go atomicIncCounter(&counter, &wg)
	go atomicIncCounter(&counter, &wg)
	wg.Wait()
	return counter
}

func ConcurrentMutexAdd() int64 {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var counter int64 = 0
	var mtx sync.Mutex
	go mutexIncCounter(&counter, &wg, &mtx)
	go mutexIncCounter(&counter, &wg, &mtx)
	wg.Wait()
	return counter
}

func main() {
	println(ConcurrentAtomicAdd())
	println(ConcurrentMutexAdd())
}
