package test

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func assertDequeModel(t *testing.T, d Deque[int], want []int) {
	t.Helper()
	if d.Len() != len(want) || d.Empty() != (len(want) == 0) || d.Cap() < len(want) || d.Fill() != (len(want) == d.Cap()) {
		t.Fatalf("invalid metadata: len=%d cap=%d empty=%v full=%v, model len=%d", d.Len(), d.Cap(), d.Empty(), d.Fill(), len(want))
	}
	if got := d.ToVec(); !slices.Equal(got, want) {
		t.Fatalf("ToVec = %v, want %v", got, want)
	}
	var each []int
	d.ForEach(func(v int) { each = append(each, v) })
	if !slices.Equal(each, want) {
		t.Fatalf("ForEach = %v, want %v", each, want)
	}
	each = each[:0]
	d.ForEachWhile(func(v int) bool { each = append(each, v); return true })
	if !slices.Equal(each, want) {
		t.Fatalf("ForEachWhile = %v, want %v", each, want)
	}
	var seq []int
	for v := range d.ToSeq() {
		seq = append(seq, v)
	}
	if !slices.Equal(seq, want) {
		t.Fatalf("ToSeq = %v, want %v", seq, want)
	}
	var indexed []int
	for i, v := range d.ToSeq2() {
		if i != len(indexed) {
			t.Fatalf("ToSeq2 index = %d, want %d", i, len(indexed))
		}
		indexed = append(indexed, v)
	}
	if !slices.Equal(indexed, want) {
		t.Fatalf("ToSeq2 = %v, want %v", indexed, want)
	}
	for i, v := range want {
		if d.Get(i) != v {
			t.Fatalf("Get(%d) = %d, want %d", i, d.Get(i), v)
		}
	}
	if len(want) == 0 {
		if d.Front().IsSome() || d.Back().IsSome() {
			t.Fatal("empty deque has endpoints")
		}
	} else if d.Front().Get() != want[0] || d.Back().Get() != want[len(want)-1] {
		t.Fatal("endpoints differ from model")
	}
}

func checkDequeOperations(t *testing.T, capacity int, ops []byte) {
	t.Helper()
	d := Deque_[int](capacity)
	var model []int
	assertDequeModel(t, d, model)
	for step, op := range ops {
		value := step*256 + int(op)
		switch op % 7 {
		case 0:
			d.PushFront(value)
			model = append([]int{value}, model...)
		case 1, 2:
			d.PushBack(value)
			model = append(model, value)
		case 3:
			got := d.PopFront()
			if len(model) == 0 {
				if got.IsSome() {
					t.Fatalf("step %d: PopFront on empty = %v", step, got)
				}
			} else {
				if got.Get() != model[0] {
					t.Fatalf("step %d: PopFront = %v, want %d", step, got, model[0])
				}
				model = model[1:]
			}
		case 4:
			got := d.PopBack()
			if len(model) == 0 {
				if got.IsSome() {
					t.Fatalf("step %d: PopBack on empty = %v", step, got)
				}
			} else {
				if got.Get() != model[len(model)-1] {
					t.Fatalf("step %d: PopBack = %v", step, got)
				}
				model = model[:len(model)-1]
			}
		case 5:
			if len(model) > 0 {
				i := int(op) % len(model)
				d.Set(i, value)
				model[i] = value
			}
		case 6:
			copy := d.ToVec()
			if len(copy) > 0 {
				copy[0] = -1
			} // ToVec must not alias storage.
		}
		assertDequeModel(t, d, model)
	}
	for len(model) > 0 {
		if got := d.PopFront().Get(); got != model[0] {
			t.Fatalf("drain = %d, want %d", got, model[0])
		}
		model = model[1:]
	}
	assertDequeModel(t, d, nil)
	capacity = d.Cap()
	d.PushFront(11)
	d.PushBack(22)
	assertDequeModel(t, d, []int{11, 22})
	if d.Cap() < capacity {
		t.Fatal("reuse unexpectedly shrank capacity")
	}
}

func TestDequeModel(t *testing.T) {
	for _, capacity := range []int{0, 1, 2, 3, 7, 16} {
		for _, seed := range []uint64{1, 42, 2026} {
			t.Run(fmt.Sprintf("cap%d/seed%d", capacity, seed), func(t *testing.T) {
				rng := rand.New(rand.NewPCG(seed, seed+1))
				ops := make([]byte, 2000)
				for i := range ops {
					ops[i] = byte(rng.Uint32())
				}
				checkDequeOperations(t, capacity, ops)
			})
		}
	}
}

func TestDequeWrappedGrowthAndEarlyStop(t *testing.T) {
	for _, front := range []bool{false, true} {
		t.Run(fmt.Sprintf("pushFront=%v", front), func(t *testing.T) {
			d := Deque_[int](4)
			for i := range 4 {
				d.PushBack(i)
			}
			d.PopFront()
			d.PopFront()
			d.PushBack(4)
			d.PushBack(5)
			want := []int{2, 3, 4, 5}
			for _, grow := range []bool{false, true} {
				if grow {
					if front {
						d.PushFront(1)
						want = []int{1, 2, 3, 4, 5}
					} else {
						d.PushBack(6)
						want = []int{2, 3, 4, 5, 6}
					}
				}
				assertDequeModel(t, d, want)
				for stop := 1; stop <= len(want); stop++ {
					var each, seq, indexed []int
					d.ForEachWhile(func(v int) bool { each = append(each, v); return len(each) < stop })
					for v := range d.ToSeq() {
						seq = append(seq, v)
						if len(seq) == stop {
							break
						}
					}
					for i, v := range d.ToSeq2() {
						if i != len(indexed) {
							t.Fatal("invalid early-stop index")
						}
						indexed = append(indexed, v)
						if len(indexed) == stop {
							break
						}
					}
					if !slices.Equal(each, want[:stop]) || !slices.Equal(seq, want[:stop]) || !slices.Equal(indexed, want[:stop]) {
						t.Fatalf("stop %d: ForEachWhile=%v ToSeq=%v ToSeq2=%v", stop, each, seq, indexed)
					}
				}
			}
		})
	}
}

func FuzzDequeModel(f *testing.F) {
	f.Add(uint8(0), []byte{0, 1, 2, 3, 4, 5, 6})
	f.Add(uint8(4), []byte{1, 1, 1, 1, 3, 3, 1, 1, 0, 4, 3})
	f.Add(uint8(1), []byte{0, 4, 1, 3, 0, 3})
	f.Fuzz(func(t *testing.T, capacity uint8, ops []byte) {
		if len(ops) > 512 {
			t.Skip()
		}
		checkDequeOperations(t, int(capacity%32), ops)
	})
}

func TestDequeJSONPreservesState(t *testing.T) {
	d := DequeOf(1, 2, 3)
	for _, bad := range []string{`[4,"bad"]`, `[`, `{}`} {
		if err := json.Unmarshal([]byte(bad), &d); err == nil {
			t.Fatalf("accepted invalid JSON %q", bad)
		}
		assertDequeModel(t, d, []int{1, 2, 3})
	}
	for _, input := range []string{`[]`, `null`, `[4,5,6]`} {
		if err := json.Unmarshal([]byte(input), &d); err != nil {
			t.Fatalf("Unmarshal(%s): %v", input, err)
		}
		var want []int
		if err := json.Unmarshal([]byte(input), &want); err != nil {
			t.Fatal(err)
		}
		assertDequeModel(t, d, want)
		data, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}
		var roundTrip Deque[int]
		if err := json.Unmarshal(data, &roundTrip); err != nil {
			t.Fatal(err)
		}
		assertDequeModel(t, roundTrip, want)
		d.PushBack(99)
		assertDequeModel(t, roundTrip, want)
	}
}

func TestDequeNilElementsAndEmptyReuse(t *testing.T) {
	for _, popFront := range []bool{false, true} {
		d := Deque_[*int](3)
		for cycle := range 20 {
			x := cycle
			d.PushFront(nil)
			d.PushBack(&x)
			if got := d.PopFront(); got.IsNone() || got.Get() != nil {
				t.Fatalf("stored nil = %v", got)
			}
			var got Opt[*int]
			if popFront {
				got = d.PopFront()
			} else {
				got = d.PopBack()
			}
			if got.IsNone() || got.Get() != &x || !d.Empty() {
				t.Fatalf("cycle %d: last pop = %v", cycle, got)
			}
			if d.PopFront().IsSome() || d.PopBack().IsSome() {
				t.Fatal("empty pop returned a value")
			}
		}
	}
}

func TestDequeInvalidIndicesPreserveState(t *testing.T) {
	for _, values := range [][]int{nil, {1}, {1, 2, 3}} {
		d := DequeOf(values...)
		for _, index := range []int{-100, -1, len(values), len(values) + 100} {
			if !panics(func() { d.Get(index) }) {
				t.Fatalf("Get(%d) did not panic for %v", index, values)
			}
			if !panics(func() { d.Set(index, 99) }) {
				t.Fatalf("Set(%d) did not panic for %v", index, values)
			}
			assertDequeModel(t, d, values)
		}
	}
}
