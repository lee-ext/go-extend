package test

import (
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestChanBasics(t *testing.T) {
	c := Chan_[int](2)
	if c.Len() != 0 || c.Cap() != 2 || !c.Empty() || c.Full() {
		t.Fatalf("new chan len=%d cap=%d", c.Len(), c.Cap())
	}
	c.Send(1)
	c.Send(2)
	if c.Len() != 2 || !c.Full() || c.Empty() {
		t.Fatalf("after send len=%d full=%v", c.Len(), c.Full())
	}
	if c.TrySend(3) {
		t.Error("TrySend on full chan should return false")
	}
	if v := c.Recv(); !v.IsSome() || v.Get() != 1 {
		t.Errorf("Recv = %v", v)
	}
	if v := c.TryRecv(); !v.IsSome() || v.Get() != 2 {
		t.Errorf("TryRecv = %v", v)
	}
	if v := c.TryRecv(); v.IsSome() {
		t.Errorf("TryRecv on drained chan = %v", v)
	}
	c.Close()
	if v := c.Recv(); v.IsSome() {
		t.Errorf("Recv on closed chan = %v", v)
	}
}

func TestChanSplit(t *testing.T) {
	c := Chan_[string](1)
	sx, rx := c.Split()
	sx.Send("hi")
	if v := rx.Recv(); !v.IsSome() || v.Get() != "hi" {
		t.Errorf("split recv = %v", v)
	}
	if !sx.TrySend("x") {
		t.Error("sender TrySend should succeed")
	}
	if sx.Len() != 1 || sx.Cap() != 1 || !sx.Full() || sx.Empty() {
		t.Errorf("sender meta len=%d cap=%d", sx.Len(), sx.Cap())
	}
	if sx.TrySend("y") {
		t.Error("sender TrySend on full chan should fail")
	}
	// drain the buffered "x" before closing, otherwise it stays readable
	if v := rx.Recv(); !v.IsSome() || v.Get() != "x" {
		t.Errorf("drain = %v", v)
	}
	sx.Close()
	if v := rx.TryRecv(); v.IsReceived() || !v.IsDisconnected() {
		t.Errorf("closed recv state = %v", v)
	}
	if rx.Len() != 0 || rx.Cap() != 1 || !rx.Empty() || rx.Full() {
		t.Errorf("receiver meta len=%d cap=%d", rx.Len(), rx.Cap())
	}
	if got := rx.Recv(); got.IsSome() {
		t.Errorf("closed Recv = %v", got)
	}
	if !rx.Iter().Empty() || rx.Iter().Len() != 0 {
		t.Error("closed chan iter should be empty")
	}
}

func TestChanForEach(t *testing.T) {
	c := Chan_[int](4)
	for i := range 3 {
		c.Send(i)
	}
	c.Close()
	var got []int
	c.ForEach(func(i int) { got = append(got, i) })
	if len(got) != 3 || got[0] != 0 || got[2] != 2 {
		t.Errorf("ForEach = %v", got)
	}

	c2 := Chan_[int](5)
	for i := range 5 {
		c2.Send(i)
	}
	c2.Close()
	count := 0
	c2.ForEachWhile(func(int) bool { count++; return count < 2 })
	if count != 2 {
		t.Errorf("ForEachWhile count = %d", count)
	}
}

func TestChanAppendSelf(t *testing.T) {
	c := Chan_[int](4)
	got := c.AppendSelf(1).AppendSelf(2)
	if got.Len() != 2 {
		t.Errorf("AppendSelf len = %d", got.Len())
	}
	if v := c.Recv(); v.Get() != 1 {
		t.Errorf("first = %v", v)
	}
	// AppendSelf returns the same channel and preserves order
	sx, rx := Chan_[int](4).Split()
	sx = sx.AppendSelf(10).AppendSelf(20)
	if a := rx.Recv(); a.Get() != 10 {
		t.Errorf("order[0] = %v", a)
	}
	if b := rx.Recv(); b.Get() != 20 {
		t.Errorf("order[1] = %v", b)
	}
	_ = sx
}

func TestChanRecvResStates(t *testing.T) {
	c := Chan_[int](1)
	_, rx := c.Split()
	if r := rx.TryRecv(); !r.IsChanEmpty() || r.IsReceived() || r.IsDisconnected() {
		t.Errorf("empty state = %v", r)
	}
	if r := rx.TryRecv(); r.ToOpt().IsSome() {
		t.Errorf("empty ToOpt = %v", r.ToOpt())
	}
	c.Send(5)
	r := rx.TryRecv()
	if !r.IsReceived() || r.ToOpt().Get() != 5 {
		t.Errorf("received state = %v", r)
	}
	c.Close()
	r2 := rx.TryRecv()
	if !r2.IsDisconnected() || r2.IsReceived() || r2.IsChanEmpty() {
		t.Errorf("disconnected state = %v", r2)
	}
	if r2.ToOpt().IsSome() {
		t.Errorf("disconnected ToOpt = %v", r2.ToOpt())
	}
}

func TestChanIter(t *testing.T) {
	c := Chan_[int](3)
	c.Send(1)
	c.Send(2)
	c.Send(3)
	c.Close()
	got := c.Iter().Collect(Vec_[int])
	if got.Len() != 3 || got[0] != 1 || got[2] != 3 {
		t.Errorf("Iter collect = %v", got)
	}
}

func TestChanConcurrency(t *testing.T) {
	c := Chan_[int](0)
	done := Promise_[int]()
	go func() {
		sum := 0
		c.ForEach(func(i int) { sum += i })
		done.Complete(sum)
	}()
	for i := 1; i <= 5; i++ {
		c.Send(i)
	}
	c.Close()
	if v := done.Await(); v != 15 {
		t.Errorf("concurrent sum = %d", v)
	}
}
