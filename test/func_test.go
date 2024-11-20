package test

import (
	"fmt"
	. "github.com/lee157953/go-extend/ext"
	"testing"
)

func TestFunc(t *testing.T) {
	vec := VecOf[int64](5, 4, 3, 2, 1)

	vec1 := Map(vec, func(t int64) float64 {
		return float64(t) / 10
	})
	fmt.Println(vec1)

	set := IntactTo(vec, Set_)
	fmt.Println(set.Contains(5))
	fmt.Println(set)

	dict := IntactTo(vec.ToSeq(), Dict_)
	fmt.Println(dict)

	vec3 := IntactTo(vec.ToRev(), Vec_)
	fmt.Println(vec3)

	dict1 := MapTo(vec, func(t int64) KV[int64, string] {
		return KV_(t, fmt.Sprintf("num:%d", t))
	}, Dict_)
	fmt.Println(dict1)

	group := VecOf(VecOf(1, 2, 3), VecOf(4, 5, 6), VecOf(7, 8, 9))
	vec4 := Flatten[int](group)
	fmt.Println(vec4)

	deque := FlattenTo[int](group, Deque_)
	fmt.Println(deque)

	group1 := VecOf("hello", " ", "world")
	vec5 := FlatMap(group1, func(t string) Vec[byte] {
		return Vec[byte](t)
	})
	fmt.Println(string(vec5))

	set1 := FlatMapTo(group1, func(t string) Vec[byte] {
		return Vec[byte](t)
	}, Set_)
	fmt.Println(set1)

	vec6 := Filter(vec, func(t int64) bool {
		return t > 2
	})
	fmt.Println(vec6)

	deque1 := FilterTo(vec, func(t int64) bool {
		return t > 2
	}, Deque_)
	fmt.Println(deque1)

	ptrVec := Map(vec,
		func(t int64) *int64 {
			if t == 3 {
				return nil
			}
			return &t
		})

	vec7 := FilterMap(ptrVec,
		func(t *int64) Opt[int64] {
			if t != nil {
				return Some(*t)
			}
			return None[int64]()
		})
	fmt.Println(vec7)

	set2 := FilterMapTo(ptrVec,
		func(t *int64) Opt[int64] {
			if t != nil {
				return Some(*t)
			}
			return None[int64]()
		}, Set_)
	fmt.Println(set2)

	sum := Reduce(vec, int64(0), func(l int64, r int64) int64 {
		return l + r
	})
	fmt.Println(sum)
}
