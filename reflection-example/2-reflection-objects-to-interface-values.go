package main

import (
	"fmt"
	"reflect"
)

type student struct {
	Name string
}

func main() {
	// FIRST EXAMPLE SHOWING CONVERT REFLECT.VALUE TO FLOAT
	floatVar := 3.14
	v := reflect.ValueOf(floatVar)
	newFloat := v.Interface().(float64)
	fmt.Println(newFloat + 1.2)

	// second example showing convert Reflect.Value to slice
	sliceVar := make([]int, 5)
	v = reflect.ValueOf(sliceVar)
	v = reflect.Append(v, reflect.ValueOf(2))
	newSlice := v.Interface().([]int)
	newSlice = append(newSlice, 4)
	fmt.Println(newSlice)

	// third example showing convert Reflect.Value to student
	stuPtr := reflect.New(reflect.TypeOf(student{}))
	stu := stuPtr.Elem()
	nameField := stu.FieldByName("Name")
	if nameField.IsValid() {
		if nameField.CanSet() {
			nameField.SetString("chong")
		}
		realStudent := stu.Interface().(student)
		fmt.Println(realStudent)
	}

}
