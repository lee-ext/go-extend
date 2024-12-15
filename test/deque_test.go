package test

import (
	"fmt"
	"github.com/lee157953/go-extend/ext"
	"testing"
)

func TestDeque(t *testing.T) {
	d := ext.Deque_[int](0)
	for i := range 5 {
		d.PushBack(i)
	}
	for i := range 5 {
		d.PushFront(i)
	}
	fmt.Println(d)
	for i, e := range d.ToSeq2() {
		fmt.Println(i, e)
	}
}
