package main

import (
	"fmt"
)

func main() {
	defer func() {
		switch r := recover().(type) {
		case nil:
			println("is nil")
		case error:
			fmt.Printf("Error: %v\n", r.Error())
		default:
			fmt.Printf("Unknown error: %#v\n", r)
		}
	}()

	for i := range 10 {
		fmt.Println(i)
	}
}

type MyErr struct {
	Msg string
}

func (e MyErr) Error() string {
	return e.Msg
}
