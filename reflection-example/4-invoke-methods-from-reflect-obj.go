package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	name string
}

func (s Student) DoHomework(number int) {
	fmt.Printf("%s is doing homework %d\n", s.name, number)
}

func main() {
	// use reflect to invoke the DoHomework of a student
	s := Student{name: "heli"}
	v := reflect.ValueOf(s)
	methodV := v.MethodByName("DoHomework")
	if methodV.IsValid() {
		in := []reflect.Value{reflect.ValueOf(55)}
		methodV.Call(in)
	}
}
