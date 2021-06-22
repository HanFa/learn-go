package main

import (
	"fmt"
	"reflect"
	"time"
)

func makeTimeFunc(f interface{}) interface{} {
	tf := reflect.TypeOf(f)
	vf := reflect.ValueOf(f)

	if tf.Kind() != reflect.Func {
		panic("expect a function")
	}
	wrapper := reflect.MakeFunc(tf, func(args []reflect.Value) (result []reflect.Value) {
		start := time.Now()
		result = vf.Call(args)
		end := time.Now()
		fmt.Printf("The function takes %v\n", end.Sub(start))
		return result
	})
	return wrapper.Interface()
}

func TimeMe() {
	time.Sleep(1 * time.Second)
}

func main() {
	timedFunc := makeTimeFunc(TimeMe).(func())
	timedFunc()
}
