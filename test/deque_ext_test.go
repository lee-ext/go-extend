package test

import (
	"encoding/json"
	"slices"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestDequeEmpty(t *testing.T) {
	d := Deque_[int](0)
	if !d.Empty() || d.Len() != 0 {
		t.Fatalf("new deque len=%d empty=%v", d.Len(), d.Empty())
	}
	if d.PopFront().IsSome() || d.PopBack().IsSome() {
		t.Error("empty Pop should be None")
	}
	if d.Front().IsSome() || d.Back().IsSome() {
		t.Error("empty Front/Back should be None")
	}
	if d.ToVec().Len() != 0 {
		t.Errorf("empty ToVec = %v", d.ToVec())
	}
	var n int
	for range d.ToSeq() {
		n++
	}
	for range d.ToSeq2() {
		n++
	}
	if n != 0 {
		t.Errorf("empty seq yielded %d", n)
	}
}

func TestDequePushPop(t *testing.T) {
	d := Deque_[int](0)
	for i := range 5 {
		d.PushBack(i)
	}
	if d.Len() != 5 || !slices.Equal(d.ToVec(), VecOf(0, 1, 2, 3, 4)) {
		t.Fatalf("PushBack = %v", d.ToVec())
	}
	if d.Front().Get() != 0 || d.Back().Get() != 4 {
		t.Errorf("Front/Back = %v/%v", d.Front(), d.Back())
	}
	for i := range 5 {
		d.PushFront(i)
	}
	if want := VecOf(4, 3, 2, 1, 0, 0, 1, 2, 3, 4); !slices.Equal(d.ToVec(), want) {
		t.Fatalf("PushFront = %v", d.ToVec())
	}
	if got := d.PopFront(); got.Get() != 4 || d.Len() != 9 {
		t.Fatalf("PopFront = %v len=%d", got, d.Len())
	}
	if got := d.PopBack(); got.Get() != 4 || d.Len() != 8 {
		t.Fatalf("PopBack = %v len=%d", got, d.Len())
	}
	for !d.Empty() {
		d.PopFront()
	}
	if d.Len() != 0 || d.PopFront().IsSome() {
		t.Errorf("drained deque = %v", d.ToVec())
	}
}

func TestDequeGrow(t *testing.T) {
	d := Deque_[int](1)
	for i := range 100 {
		d.PushBack(i)
	}
	if d.Len() != 100 {
		t.Fatalf("len = %d", d.Len())
	}
	for i, e := range d.ToVec() {
		if e != i {
			t.Fatalf("ToVec[%d] = %d", i, e)
		}
	}
	e := Deque_[int](1)
	for i := range 100 {
		e.PushFront(i)
	}
	if want := VecInit(100, func(i int) int { return 99 - i }); !slices.Equal(e.ToVec(), want) {
		t.Fatalf("PushFront grow = %v", e.ToVec()[:5])
	}
	if d.Cap() < 100 || e.Cap() < 100 {
		t.Errorf("cap d=%d e=%d", d.Cap(), e.Cap())
	}
}

func TestDequeWrapAround(t *testing.T) {
	d := DequeOf(1, 2, 3)
	d.PushFront(0)
	if want := VecOf(0, 1, 2, 3); !slices.Equal(d.ToVec(), want) {
		t.Fatalf("wrap ToVec = %v", d.ToVec())
	}
	if d.Len() != 4 {
		t.Fatalf("wrap len = %d", d.Len())
	}
	for i, want := range []int{0, 1, 2, 3} {
		if got := d.Get(i); got != want {
			t.Errorf("Get(%d) = %d", i, got)
		}
	}
	d.Set(0, 9)
	if d.Get(0) != 9 || d.Front().Get() != 9 {
		t.Errorf("Set(0) = %d", d.Get(0))
	}
	var seq []int
	for i, e := range d.ToSeq2() {
		if i != len(seq) {
			t.Fatalf("ToSeq2 index = %d, want %d", i, len(seq))
		}
		seq = append(seq, e)
	}
	if !slices.Equal(seq, []int{9, 1, 2, 3}) {
		t.Errorf("ToSeq2 = %v", seq)
	}
	count := 0
	d.ForEach(func(int) { count++ })
	if count != 4 {
		t.Errorf("ForEach count = %d", count)
	}
	count = 0
	d.ForEachWhile(func(int) bool { count++; return count < 2 })
	if count != 2 {
		t.Errorf("ForEachWhile count = %d", count)
	}
}

func TestDequeIndexPanic(t *testing.T) {
	d := DequeOf(1, 2, 3)
	for _, bad := range []int{-1, 3, 99} {
		if !panics(func() { d.Get(bad) }) {
			t.Errorf("Get(%d) should panic", bad)
		}
	}
	if !panics(func() { d.Set(3, 0) }) {
		t.Error("Set(3) should panic")
	}
}

func TestDequeMisc(t *testing.T) {
	d := Deque_[int](0)
	if got := d.AppendSelf(1).AppendSelf(2); got.Len() != 2 || got.Get(1) != 2 {
		t.Errorf("AppendSelf = %v", got)
	}
	if s := DequeOf(1, 2).String(); s != "deque[1 2]" {
		t.Errorf("String = %q", s)
	}
	if got := DequeOf(1, 2).Iter().Collect(Vec_[int]); !slices.Equal(got, VecOf(1, 2)) {
		t.Errorf("Iter collect = %v", got)
	}
	b, err := json.Marshal(DequeOf(1, 2, 3))
	if err != nil || string(b) != "[1,2,3]" {
		t.Errorf("Marshal = %s, %v", b, err)
	}
	var d2 Deque[int]
	if err := json.Unmarshal([]byte(`[4,5]`), &d2); err != nil || !slices.Equal(d2.ToVec(), VecOf(4, 5)) {
		t.Errorf("Unmarshal = %v, %v", d2, err)
	}
}

func TestDequeOfOrder(t *testing.T) {
	d := DequeOf("a", "b", "c")
	if d.Len() != 3 || d.Get(0) != "a" || d.Get(2) != "c" {
		t.Fatalf("DequeOf = %v", d.ToVec())
	}
	if d.PopBack().Get() != "c" || d.PopFront().Get() != "a" {
		t.Errorf("pops = %v", d.ToVec())
	}
	if d.Len() != 1 || d.Get(0) != "b" {
		t.Errorf("rest = %v", d.ToVec())
	}
}
