package main

import (
	"fmt"
	"go-extend/ext"
	"time"
)

func main() {
	println("hello")
	a := ext.Actor_(4, func(a any) {
		fmt.Printf("%#v\n", a)
	})
	a.Launch(func() {
		panic("a")
	})
	a.Launch(func() {
		panic("b")
	})
	time.Sleep(time.Second)
}
