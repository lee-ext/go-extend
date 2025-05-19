package test

import (
	"fmt"
	"github.com/lee-ext/go-extend/ext"
	"testing"
	"time"
)

func TestActor(t *testing.T) {
	actor := ext.Actor_(8, ext.DeferFn_)
	defer actor.Close()
	for i := range 5 {
		actor.Launch(func() {
			if i == 3 {
				panic(fmt.Errorf("task: %v", i))
			}
			println(fmt.Printf("task: %v\n", i))
		})
	}
	p := ext.Promise_[string]()
	actor.Launch(func() {
		p.Complete("hello actor")
	})
	println(p.Await().Get())
	time.Sleep(time.Second)
}
