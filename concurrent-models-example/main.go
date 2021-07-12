package main

import (
	"fmt"
	"io"
	"learn-go-concurrent/pool"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// DBConnection	定义的一个资源
type DBConnection struct {
	id int32
}

func (D DBConnection) Close() error {
	fmt.Println("database closed , #" + fmt.Sprint(D.id))
	return nil
}

var counter int32

func Factory() (io.Closer, error) {
	atomic.AddInt32(&counter, 1)
	return DBConnection{
		id: counter,
	}, nil
}

func performQuery(query int, pool *pool.Pool) {
	defer wg.Done()

	resource, err := pool.AcquireResource()
	if err != nil {
		log.Fatalln(err)
	}
	defer pool.ReleaseResource(resource)

	t := rand.Int()%10 + 1
	time.Sleep(time.Duration(t) * time.Second)
	fmt.Println("finish query" + fmt.Sprint(query))

}

var wg sync.WaitGroup

func main() {
	p, err := pool.New(Factory, 5)
	if err != nil {
		log.Fatalln(err)
	}

	num := 10
	wg.Add(num)
	for id := 0; id < num; id++ {
		go performQuery(id, p)
	}
	wg.Wait()

	p.Close()

	fmt.Println("pool model run ends")
}

/*
// runner model demo
func createTask() func(int) {
	return func(id int) {
		time.Sleep(1 * time.Second)
		fmt.Printf("Task complete #%d \n", id)
	}
}

func main() {

	r := runner.New(4 * time.Second)

	r.AddTask(createTask(), createTask(), createTask())

	err := r.Start()

	switch err {
	case runner.ErrInterrupt:
		fmt.Println("tasks interrupted")
	case runner.ErrTimeout:
		fmt.Println("tasks timeout")
	default:
		fmt.Println("all tasks finished")
	}
}
*/
