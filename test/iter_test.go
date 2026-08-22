package test

import (
	"fmt"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestIterFunc(t *testing.T) {
	s := VecOf(VecOf(0, 1, 2, 3, 4, 5),
		VecOf(5, 6, 7, 8, 9)).Iter().
		//FlatMap(func(v Vec[int]) Vec[int] {
		//	println("FlatMap")
		//	return v
		//}).
		Flatten[Vec[int]]().
		FilterMap(func(i int) Opt[int] {
			println("FilterMap")
			if i > 0 {
				return Some(i + 1)
			}
			return None[int]()
		}).
		Filter(func(i int) bool {
			println("Filter")
			return i%2 == 0
		}).Map(func(i int) KV[int, string] {
		println("Map")
		return KV_(i, fmt.Sprintf("val~%d", i))
	}).Collect(MDict_[int, string])
	fmt.Printf("%v\n", s)
}
