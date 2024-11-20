package test

import (
	"fmt"
	"github.com/lee157953/go-extend/ext"
	"testing"
	"time"
)

func TestActor(t *testing.T) {
	actor := ext.Actor_(8, ext.DefaultDeferFn)
	defer actor.Close()
	for i := range 5 {
		actor.Launch(func() {
			if i == 3 {
				panic(fmt.Errorf("task: %v", i))
			}
			println(fmt.Printf("task: %v\n", i))
		})
	}
	p, fn := ext.Promise_[string](0)
	actor.Launch(func() {
		fn("hello actor")
	})
	println(p.Await().Get())
	time.Sleep(time.Second)
}
