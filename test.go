package main

import (
	"fmt"
	"reflect"
)

type test struct {
	name string `json:"product_name" sql:"product_name" product:"id"`
	sex  bool
}

func(t *test)GetName()string {
	return t.name
}

func main() {
	a := test{name:"xxz", sex: true}
	fmt.Println(reflect.ValueOf(&a).Elem())
}
