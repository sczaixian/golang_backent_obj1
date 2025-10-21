package main

import (
	"fmt"
	"reflect"
)

type MyStruct struct {
	Name string
}

func (m *MyStruct) SayHello() {
	fmt.Printf("Hello, I'm %s\n", m.Name)
}

func (m *MyStruct) SayGoodbye() {
	fmt.Printf("Goodbye, %s\n", m.Name)
}

func main() {
	// 创建结构体实例
	my := &MyStruct{Name: "John"}

	// 获取反射值对象
	value := reflect.ValueOf(my)

	// 通过值反射对象获取方法（注意：方法必须是对外暴露的，即首字母大写）
	method := value.MethodByName("SayHello")
	if method.IsValid() {
		// 调用方法
		method.Call(nil)
	} else {
		fmt.Println("Method not found.")
	}
	fmt.Println("-------------------------------------------------")
	// 另一种方式：通过类型反射对象获取方法信息
	typ := reflect.TypeOf(my)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		fmt.Printf("Method %d: %s\n", i, m.Name)
	}
}
