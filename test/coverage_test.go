package test

import (
	"cmp"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/lee-ext/go-extend/ext"
)

// CmpVec covers the Func-suffixed and mutating methods not exercised elsewhere.
func TestCmpVecFuncs(t *testing.T) {
	v := CmpVecOf(3, 1, 2)
	var seen []int
	v.ForEach(func(i int) { seen = append(seen, i) })
	if !slices.Equal(seen, []int{3, 1, 2}) {
		t.Errorf("ForEach = %v", seen)
	}
	n := 0
	v.ForEachWhile(func(int) bool { n++; return n < 2 })
	if n != 2 {
		t.Errorf("ForEachWhile n = %d", n)
	}
	if got := v.Iter().Collect(CmpVec_[int]); !slices.Equal(got, v) {
		t.Errorf("Iter collect = %v", got)
	}
	if v.First().Get() != 3 || v.Last().Get() != 2 {
		t.Error("bad First/Last")
	}
	v.SortFunc(cmp.Compare[int])
	if !v.IsSortedFunc(cmp.Compare[int]) || !slices.Equal(v, CmpVec[int]{1, 2, 3}) {
		t.Errorf("SortFunc = %v", v)
	}
	stable := CmpVecOf(T2_(1, "a"), T2_(1, "b"), T2_(0, "c"))
	stable.SortStableFunc(func(a, b T2[int, string]) int { return cmp.Compare(a.V0, b.V0) })
	if stable[1].V1 != "a" || stable[2].V1 != "b" {
		t.Errorf("SortStableFunc = %v", stable)
	}
	if v.MinFunc(cmp.Compare[int]) != 1 || v.MaxFunc(cmp.Compare[int]) != 3 {
		t.Error("bad MinFunc/MaxFunc")
	}
	if idx, ok := v.BinarySearchFunc(2, cmp.Compare[int]); !ok || idx != 1 {
		t.Errorf("BinarySearchFunc = %d,%v", idx, ok)
	}
	if v.IndexFunc(func(i int) bool { return i == 3 }) != 2 {
		t.Error("bad IndexFunc")
	}
	if !v.ContainsFunc(func(i int) bool { return i == 1 }) || v.ContainsFunc(func(i int) bool { return i == 9 }) {
		t.Error("bad ContainsFunc")
	}
}

func TestCmpVecMutate(t *testing.T) {
	v := CmpVec_[int](0)
	v.Append(1)
	v.Appends(2, 3)
	if !slices.Equal(v, CmpVec[int]{1, 2, 3}) {
		t.Fatalf("Append(s) = %v", v)
	}
	if got := v.AppendSelf(4); !slices.Equal(got, CmpVec[int]{1, 2, 3, 4}) {
		t.Errorf("AppendSelf = %v", got)
	}
	v.Insert(0, 0)
	v.Replace(1, 3, 9)
	v.RemoveRange(0, 1)
	v.RemoveAt(0)
	if !slices.Equal(v, CmpVec[int]{3, 4}) {
		t.Fatalf("mutate = %v", v)
	}
	v.Repeat(2)
	if !slices.Equal(v, CmpVec[int]{3, 4, 3, 4}) {
		t.Fatalf("Repeat = %v", v)
	}
	v.Grow(8)
	if v.Cap() < v.Len()+8 {
		t.Errorf("Grow cap = %d", v.Cap())
	}
	v.Clip()
	if v.Cap() != v.Len() {
		t.Errorf("Clip cap = %d", v.Cap())
	}
	cp := v.Clone()
	cp.Append(99)
	if v.Len() != 4 {
		t.Errorf("Clone aliased source: %v", v)
	}
	v.Reverse()
	if !slices.Equal(v, CmpVec[int]{4, 3, 4, 3}) {
		t.Errorf("Reverse = %v", v)
	}
	dup := CmpVecOf(1, 1, 2)
	dup.CompactFunc(func(a, b int) bool { return a == b })
	if !slices.Equal(dup, CmpVec[int]{1, 2}) {
		t.Errorf("CompactFunc = %v", dup)
	}
	sh := CmpVecInit(40, func(i int) int { return i })
	sh.Shuffle()
	sorted := sh.Clone()
	sorted.Sort()
	if !slices.Equal(sorted, CmpVecInit(40, func(i int) int { return i })) {
		t.Error("Shuffle must keep a permutation")
	}
	v.Clear()
	if !v.Empty() || v.Len() != 0 {
		t.Errorf("Clear = %v", v)
	}
	e := CmpVec[int]{}
	if e.First().IsSome() || e.Last().IsSome() || e.Pop().IsSome() {
		t.Error("empty CmpVec accessors should be None")
	}
}

func TestReceiverForEach(t *testing.T) {
	c := Chan_[int](3)
	for i := range 3 {
		c.Send(i)
	}
	c.Close()
	_, rx := c.Split()
	var got []int
	rx.ForEach(func(i int) { got = append(got, i) })
	if !slices.Equal(got, []int{0, 1, 2}) {
		t.Errorf("Receiver.ForEach = %v", got)
	}

	c2 := Chan_[int](5)
	for i := range 5 {
		c2.Send(i)
	}
	c2.Close()
	_, rx2 := c2.Split()
	count := 0
	rx2.ForEachWhile(func(int) bool { count++; return count < 2 })
	if count != 2 {
		t.Errorf("Receiver.ForEachWhile count = %d", count)
	}
}

func TestMapWhile(t *testing.T) {
	src := VecOf(1, 2, 3, 0, 4)
	// stops at the first None
	got := MapWhile(src, func(i int) Opt[int] {
		if i == 0 {
			return None[int]()
		}
		return Some(i * 10)
	})
	if !slices.Equal(got, VecOf(10, 20, 30)) {
		t.Errorf("MapWhile = %v", got)
	}
	got2 := MapWhileTo(src, func(i int) Opt[int] {
		if i == 0 {
			return None[int]()
		}
		return Some(i)
	}, Deque_[int])
	if !slices.Equal(got2.ToVec(), VecOf(1, 2, 3)) {
		t.Errorf("MapWhileTo = %v", got2)
	}
	// all Some -> no early stop
	all := MapWhile(VecOf(1, 2), func(i int) Opt[int] { return Some(i) })
	if !slices.Equal(all, VecOf(1, 2)) {
		t.Errorf("MapWhile all = %v", all)
	}
}

func TestFilterLenBranches(t *testing.T) {
	// filterLen has three size branches: <8, [8,32), >=32
	for _, size := range []int{3, 16, 40} {
		src := VecInit(size, func(i int) int { return i })
		got := Filter(src, func(i int) bool { return i%2 == 0 })
		want := size/2 + size%2
		if got.Len() != want {
			t.Errorf("size %d Filter len = %d, want %d", size, got.Len(), want)
		}
		fm := FilterMap(src, func(i int) Opt[int] {
			if i%2 == 0 {
				return Some(i)
			}
			return None[int]()
		})
		if fm.Len() != want {
			t.Errorf("size %d FilterMap len = %d", size, fm.Len())
		}
	}
}

func TestIterEmptyBranches(t *testing.T) {
	src := VecOf(1, 2, 3)
	empty := Vec[int]{}
	if src.Iter().Map(func(i int) int { return i }).Empty() {
		t.Error("Map non-empty should not be Empty")
	}
	if !empty.Iter().Map(func(i int) int { return i }).Empty() {
		t.Error("Map empty should be Empty")
	}
	if src.Iter().Filter(func(int) bool { return true }).Empty() {
		t.Error("Filter non-empty should not be Empty")
	}
	if !empty.Iter().Filter(func(int) bool { return true }).Empty() {
		t.Error("Filter empty should be Empty")
	}
	if src.Iter().FilterMap(func(i int) Opt[int] { return Some(i) }).Empty() {
		t.Error("FilterMap non-empty should not be Empty")
	}
	if !empty.Iter().FilterMap(func(i int) Opt[int] { return Some(i) }).Empty() {
		t.Error("FilterMap empty should be Empty")
	}
	flat := src.Iter().FlatMap(func(i int) Vec[int] { return VecOf(i) })
	if flat.Empty() || flat.Len() != 3 {
		t.Errorf("FlatMap meta len=%d empty=%v", flat.Len(), flat.Empty())
	}
	// ForEachWhile early-exit on filter / filterMap / flatMap / flatten
	count := 0
	src.Iter().Filter(func(int) bool { return true }).ForEachWhile(func(int) bool { count++; return count < 2 })
	if count != 2 {
		t.Errorf("Filter ForEachWhile = %d", count)
	}
	count = 0
	src.Iter().FilterMap(func(i int) Opt[int] { return Some(i) }).ForEachWhile(func(int) bool { count++; return count < 2 })
	if count != 2 {
		t.Errorf("FilterMap ForEachWhile = %d", count)
	}
	count = 0
	src.Iter().FlatMap(func(i int) Vec[int] { return VecOf(i, i) }).ForEachWhile(func(int) bool { count++; return count < 3 })
	if count != 3 {
		t.Errorf("FlatMap ForEachWhile = %d", count)
	}
	count = 0
	VecOf(VecOf(1, 2), VecOf(3)).Iter().Flatten[Vec[int]]().ForEachWhile(func(int) bool { count++; return count < 2 })
	if count != 2 {
		t.Errorf("Flatten ForEachWhile = %d", count)
	}
}

func TestDequeToSeqWrap(t *testing.T) {
	d := DequeOf(1, 2, 3)
	d.PushFront(0) // forces head > tail (wrap-around)
	var seq []int
	for e := range d.ToSeq() {
		seq = append(seq, e)
	}
	if !slices.Equal(seq, []int{0, 1, 2, 3}) {
		t.Errorf("ToSeq wrap = %v", seq)
	}
	var idx []int
	for i := range d.ToSeq2() {
		idx = append(idx, i)
	}
	if !slices.Equal(idx, []int{0, 1, 2, 3}) {
		t.Errorf("ToSeq2 wrap indices = %v", idx)
	}
	// ToSeq/ToSeq2 must stop early when yield returns false
	first := 0
	for e := range d.ToSeq() {
		first = e
		break
	}
	if first != 0 {
		t.Errorf("ToSeq first = %d", first)
	}
}

func TestDequeForEachWhileWrap(t *testing.T) {
	d := DequeOf(1, 2, 3)
	d.PushFront(0) // head > tail
	count := 0
	d.ForEachWhile(func(int) bool { count++; return count < 2 })
	if count != 2 {
		t.Errorf("ForEachWhile wrap count = %d", count)
	}
	// exhaust the wrap branch fully
	total := 0
	d.ForEach(func(int) { total++ })
	if total != 4 {
		t.Errorf("ForEach wrap total = %d", total)
	}
}

func TestBytesWriteUInt(t *testing.T) {
	b := Bytes_(8)
	b.WriteUInt8(0, 255)
	b.WriteUInt32(1, 0x01020304)
	if b.ReadUInt8(0) != 255 {
		t.Errorf("WriteUInt8 = %d", b.ReadUInt8(0))
	}
	if b.ReadUInt32(1) != 0x01020304 {
		t.Errorf("WriteUInt32 = %#x", b.ReadUInt32(1))
	}
}

func TestBytes2BitMapValue(t *testing.T) {
	raw := make([]byte, 1)
	bm := Bytes2BitMap_(raw)
	bm.Set(0, 3)
	if len(bm.Value()) != 1 {
		t.Errorf("Value len = %d", len(bm.Value()))
	}
	if bm.Value()[0]&0b11 != 0b11 {
		t.Errorf("Value bits = %08b", bm.Value()[0])
	}
}

func TestOptFGetElseSome(t *testing.T) {
	o := OptF_[float64](2.5, true)
	if o.GetElse(func() float64 { return -1 }) != 2.5 {
		t.Errorf("OptF GetElse some = %v", o.GetElse(func() float64 { return -1 }))
	}
	n := OptF_[float64](0, false)
	if n.GetElse(func() float64 { return -1 }) != -1 {
		t.Errorf("OptF GetElse none = %v", n.GetElse(func() float64 { return -1 }))
	}
}

func TestOptNzGetOrGetElse(t *testing.T) {
	some := NzOpt_(5)
	if some.GetOr(1) != 5 || some.GetElse(func() int { return 2 }) != 5 {
		t.Errorf("NzOpt some fallbacks = %v", some)
	}
	none := NzOpt_(0)
	if none.GetOr(7) != 7 || none.GetElse(func() int { return 8 }) != 8 {
		t.Errorf("NzOpt none fallbacks = %v", none)
	}
}

func TestLaunchSturdyRestarts(t *testing.T) {
	// LaunchSturdy re-runs fn after each return/panic; count restarts atomically.
	var count atomic.Int64
	LaunchSturdy(func() {
		count.Add(1)
	}, func(any) {}, 5*time.Millisecond)
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if count.Load() >= 3 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if count.Load() < 3 {
		t.Errorf("LaunchSturdy should restart, count = %d", count.Load())
	}
}

func TestLaunchSturdyRecoversPanic(t *testing.T) {
	var recovered atomic.Int64
	var count atomic.Int64
	done := Promise_[Unit]()
	LaunchSturdy(func() {
		if count.Add(1) == 1 {
			panic("boom")
		}
		if !done.Completed() {
			done.Complete(Unit{})
		}
	}, func(r any) {
		if r != nil {
			recovered.Add(1)
		}
	}, 5*time.Millisecond)
	done.Await()
	if recovered.Load() < 1 {
		t.Errorf("deferFn should observe panic, recovered = %d", recovered.Load())
	}
}
