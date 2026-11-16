package test

import (
	"maps"
	"slices"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestIterMapFilter(t *testing.T) {
	src := VecOf(1, 2, 3, 4, 5)
	got := src.Iter().Map(func(i int) int { return i * 2 }).Collect(Vec_[int])
	if !slices.Equal(got, VecOf(2, 4, 6, 8, 10)) {
		t.Errorf("Map = %v", got)
	}
	got = src.Iter().Filter(func(i int) bool { return i%2 == 0 }).Collect(Vec_[int])
	if !slices.Equal(got, VecOf(2, 4)) {
		t.Errorf("Filter = %v", got)
	}
	got = src.Iter().
		Filter(func(i int) bool { return i > 1 }).
		Map(func(i int) int { return i * 10 }).
		Collect(Vec_[int])
	if !slices.Equal(got, VecOf(20, 30, 40, 50)) {
		t.Errorf("chained = %v", got)
	}
}

func TestIterFilterMap(t *testing.T) {
	src := VecOf(1, 2, 3, 4)
	got := src.Iter().
		FilterMap(func(i int) Opt[int] {
			if i%2 == 0 {
				return Some(i)
			}
			return None[int]()
		}).
		Collect(Vec_[int])
	if !slices.Equal(got, VecOf(2, 4)) {
		t.Errorf("FilterMap = %v", got)
	}
}

func TestIterFlatMap(t *testing.T) {
	src := VecOf(1, 2, 3)
	got := src.Iter().
		FlatMap(func(i int) Vec[int] { return VecOf(i, i*10) }).
		Collect(Vec_[int])
	if !slices.Equal(got, VecOf(1, 10, 2, 20, 3, 30)) {
		t.Errorf("FlatMap = %v", got)
	}
	// FlatMapEager produces the same order but computes total length eagerly
	got = src.Iter().
		FlatMapEager(func(i int) Vec[int] { return VecOf(i, i*10) }).
		Collect(Vec_[int])
	if !slices.Equal(got, VecOf(1, 10, 2, 20, 3, 30)) {
		t.Errorf("FlatMapEager = %v", got)
	}
	eager := src.Iter().FlatMapEager(func(i int) Vec[int] { return VecOf(i, i*10) })
	if eager.Len() != 6 || eager.Empty() {
		t.Errorf("FlatMapEager meta len=%d", eager.Len())
	}
}

func TestIterFlatten(t *testing.T) {
	group := VecOf(VecOf(1, 2), VecOf(3), VecOf(4, 5, 6))
	got := group.Iter().Flatten[Vec[int]]().Collect(Vec_[int])
	if !slices.Equal(got, VecOf(1, 2, 3, 4, 5, 6)) {
		t.Errorf("Flatten = %v", got)
	}
	// Flatten on a non-iterable element type must panic
	if !panics(func() {
		VecOf(1, 2, 3).Iter().Flatten[Vec[int]]().Collect(Vec_[int])
	}) {
		t.Error("Flatten on flat Vec should panic")
	}
}

func TestIterReduce(t *testing.T) {
	src := VecOf(1, 2, 3, 4)
	if sum := src.Iter().Reduce(0, func(acc, i int) int { return acc + i }); sum != 10 {
		t.Errorf("Reduce sum = %d", sum)
	}
	if prod := src.Iter().Reduce(1, func(acc, i int) int { return acc * i }); prod != 24 {
		t.Errorf("Reduce product = %d", prod)
	}
	empty := Vec[int]{}
	if s := empty.Iter().Reduce(7, func(acc, i int) int { return acc + i }); s != 7 {
		t.Errorf("Reduce on empty = %d", s)
	}
}

func TestIterLazy(t *testing.T) {
	// Map/Filter must not run until a terminal op consumes the iterator.
	calls := 0
	src := VecOf(1, 2, 3)
	lazy := src.Iter().Map(func(i int) int { calls++; return i })
	if calls != 0 {
		t.Errorf("Map should be lazy, calls=%d", calls)
	}
	lazy.Collect(Vec_[int])
	if calls != 3 {
		t.Errorf("after Collect calls=%d", calls)
	}
}

func TestIterMeta(t *testing.T) {
	src := VecOf(1, 2, 3)
	if it := src.Iter(); it.Len() != 3 || it.Empty() {
		t.Errorf("Iter meta len=%d", it.Len())
	}
	empty := Vec[int]{}
	if it := empty.Iter(); it.Len() != 0 || !it.Empty() {
		t.Errorf("empty Iter meta len=%d", it.Len())
	}
	// Filter/FilterMap report an estimated length, not the exact one
	if it := src.Iter().Filter(func(int) bool { return true }); it.Empty() {
		t.Error("filtered non-empty iter should not be Empty")
	}
}

func TestIterForEachWhile(t *testing.T) {
	src := VecOf(1, 2, 3, 4, 5)
	count := 0
	src.Iter().Map(func(i int) int { return i }).ForEachWhile(func(int) bool {
		count++
		return count < 3
	})
	if count != 3 {
		t.Errorf("ForEachWhile count = %d", count)
	}
	count = 0
	src.Iter().Filter(func(i int) bool { return i%2 == 1 }).ForEachWhile(func(int) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("Filter ForEachWhile count = %d", count)
	}
	count = 0
	src.Iter().FilterMap(func(i int) Opt[int] { return Some(i) }).ForEachWhile(func(int) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("FilterMap ForEachWhile count = %d", count)
	}
}

func TestIterCollectTargets(t *testing.T) {
	src := VecOf(1, 2, 2, 3)
	if s := src.Iter().Collect(Set_[int]); s.Len() != 3 {
		t.Errorf("Collect Set = %v", s)
	}
	if d := src.Iter().Map(KeyOf(func(i int) int { return i })).Collect(Dict_[int, int]); d.Len() != 3 {
		t.Errorf("Collect Dict = %v", d)
	}
	if dq := src.Iter().Collect(Deque_[int]); dq.Len() != 4 {
		t.Errorf("Collect Deque = %v", dq)
	}
}

func TestKeyOf(t *testing.T) {
	fn := KeyOf(func(s string) int { return len(s) })
	if kv := fn("abc"); kv.K != 3 || kv.V != "abc" {
		t.Errorf("KeyOf = %v", kv)
	}
	got := VecOf("a", "bb", "ccc").Iter().Map(KeyOf(func(s string) int { return len(s) })).
		Collect(Dict_[int, string])
	if !maps.Equal(got, Dict[int, string]{1: "a", 2: "bb", 3: "ccc"}) {
		t.Errorf("KeyOf dict = %v", got)
	}
}

func TestIterFromDictAndSet(t *testing.T) {
	d := DictOf(KV_("a", 1), KV_("b", 2))
	sum := d.Iter().Reduce(0, func(acc int, kv KV[string, int]) int { return acc + kv.V })
	if sum != 3 {
		t.Errorf("Dict Iter reduce = %d", sum)
	}
	s := SetOf(1, 2, 3)
	if v := s.Iter().Collect(Vec_[int]); v.Len() != 3 {
		t.Errorf("Set Iter collect = %v", v)
	}
	if d.Iter().Len() != 2 || d.Iter().Empty() {
		t.Errorf("Dict Iter meta len=%d", d.Iter().Len())
	}
}
