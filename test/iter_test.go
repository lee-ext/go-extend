package test

import (
	"fmt"
	"maps"
	"slices"
	"testing"
	"unsafe"

	. "github.com/lee-ext/go-extend/ext"
)

func TestIterStorageSize(t *testing.T) {
	// Source validation must not add per-instance metadata to every iterator.
	if got, want := unsafe.Sizeof(Iter[int]{}), unsafe.Sizeof(Iterator[int](nil)); got != want {
		t.Fatalf("Iter size = %d, want source interface size %d", got, want)
	}
}

func TestIterPipeline(t *testing.T) {
	got := VecOf(VecOf(0, 1, 2, 3, 4, 5), VecOf(5, 6, 7, 8, 9)).Iter().
		FlatMap(func(v Vec[int]) Vec[int] { return v }).
		FilterMap(func(i int) Opt[int] {
			if i > 0 {
				return Some(i + 1)
			}
			return None[int]()
		}).
		Filter(func(i int) bool { return i%2 == 0 }).
		Map(KeyOf(func(i int) int { return i })).Collect(MDict_[int, int])
	want := MDict[int, int]{2: {2}, 4: {4}, 6: {6, 6}, 8: {8}, 10: {10}}
	if !maps.EqualFunc(got, want, func(a, b Vec[int]) bool { return slices.Equal(a, b) }) {
		t.Fatalf("pipeline = %v, want %v", got, want)
	}
}

func TestIterShortCircuit(t *testing.T) {
	for _, kind := range []string{"map", "filter", "filterMap", "flatMap", "flatten", "eager"} {
		for stop := 1; stop <= 4; stop++ {
			t.Run(fmt.Sprintf("%s/stop%d", kind, stop), func(t *testing.T) {
				calls := 0
				base := VecOf(0, 1, 2, 3, 4, 5, 6, 7).Iter()
				var it Iter[int]
				switch kind {
				case "map":
					it = VecOf(0, 2, 4, 6).Iter().Map(func(v int) int { calls++; return v })
				case "filter":
					it = base.Filter(func(v int) bool { calls++; return v%2 == 0 })
				case "filterMap":
					it = base.FilterMap(func(v int) Opt[int] { calls++; return Opt_(v, v%2 == 0) })
				case "flatMap":
					it = VecOf(0, 4, 8).Iter().FlatMap(func(v int) Vec[int] { calls++; return VecOf(v, v+2) })
				case "flatten":
					it = VecOf(Vec[int]{}, VecOf(0, 2), Vec[int]{}, VecOf(4, 6), VecOf(8)).Iter().Flatten[Vec[int]]()
				case "eager":
					it = VecOf(0, 4, 8).Iter().FlatMapEager(func(v int) Vec[int] { calls++; return VecOf(v, v+2) })
				}
				var got []int
				it.ForEachWhile(func(v int) bool { got = append(got, v); return len(got) < stop })
				want := []int{0, 2, 4, 6}[:stop]
				if !slices.Equal(got, want) {
					t.Fatalf("yielded %v, want %v", got, want)
				}
				wantCalls := map[string]int{"map": stop, "filter": 2*stop - 1, "filterMap": 2*stop - 1, "flatMap": (stop + 1) / 2, "flatten": 0, "eager": 3}[kind]
				if calls != wantCalls {
					t.Fatalf("callback calls = %d, want %d", calls, wantCalls)
				}
			})
		}
	}
}

func TestIterRepeatabilityAndLaziness(t *testing.T) {
	calls := 0
	it := VecOf(1, 2, 3).Iter().Map(func(v int) int { calls++; return v * 2 }).
		Filter(func(v int) bool { return v > 2 }).
		FilterMap(func(v int) Opt[string] { return Some(fmt.Sprint(v)) })
	// Metadata is a source-based estimate, and must not evaluate lazy callbacks.
	_ = it.Len()
	_ = it.Empty()
	if calls != 0 {
		t.Fatal("construction/metadata evaluated lazy pipeline")
	}
	for pass := 1; pass <= 3; pass++ {
		got := it.Collect(Vec_[string])
		if !slices.Equal(got, []string{"4", "6"}) || calls != pass*3 {
			t.Fatalf("pass %d: got %v, calls %d", pass, got, calls)
		}
		got[0] = "changed" // Collections must have independent storage.
	}
}

func TestIterEmptyPipelines(t *testing.T) {
	unexpected := func(int) int { t.Fatal("callback on empty source"); return 0 }
	empty := Vec[int]{}.Iter()
	for name, it := range map[string]Iter[int]{
		"map":           empty.Map(unexpected),
		"filter":        empty.Filter(func(v int) bool { unexpected(v); return true }),
		"filterMap":     empty.FilterMap(func(v int) Opt[int] { return Some(unexpected(v)) }),
		"flatMap":       empty.FlatMap(func(v int) Vec[int] { return VecOf(unexpected(v)) }),
		"eager":         empty.FlatMapEager(func(v int) Vec[int] { return VecOf(unexpected(v)) }),
		"flatten":       VecOf(Vec[int](nil), Vec[int]{}).Iter().Flatten[Vec[int]](),
		"allFiltered":   VecOf(1, 2).Iter().Filter(func(int) bool { return false }),
		"allNone":       VecOf(1, 2).Iter().FilterMap(func(int) Opt[int] { return None[int]() }),
		"allInnerEmpty": VecOf(1, 2).Iter().FlatMap(func(int) Vec[int] { return nil }),
	} {
		t.Run(name, func(t *testing.T) {
			if got := it.Collect(Vec_[int]); len(got) != 0 {
				t.Fatalf("Collect = %v", got)
			}
			it.ForEach(func(v int) { unexpected(v) })
			it.ForEachWhile(func(v int) bool { unexpected(v); return false })
			if got := it.Reduce(42, func(acc, v int) int { return unexpected(v) }); got != 42 {
				t.Fatalf("empty Reduce = %d, want seed 42", got)
			}
		})
	}
}

func TestIterChannelConsumption(t *testing.T) {
	ch := Chan_[int](4)
	for _, v := range []int{1, 2, 3, 4} {
		ch.Send(v)
	}
	ch.Close()
	var got []int
	ch.Iter().Filter(func(v int) bool { return v%2 == 0 }).ForEachWhile(func(v int) bool {
		got = append(got, v)
		return false
	})
	if !slices.Equal(got, []int{2}) {
		t.Fatalf("first traversal = %v", got)
	}
	if rest := ch.Iter().Collect(Vec_[int]); !slices.Equal(rest, []int{3, 4}) {
		t.Fatalf("early stop consumed extra channel values: %v", rest)
	}
}

func TestIterFlattenRetainsSource(t *testing.T) {
	groups := VecOf(VecOf(1, 2), VecOf(3))
	flat := groups.Iter().Flatten[Vec[int]]()
	// Replacing an inner slice of equal length remains visible: Flatten must not
	// copy the outer collection merely to support consuming channel sources.
	groups[0] = VecOf(4, 5)
	if got := flat.Collect(Vec_[int]); !slices.Equal(got, []int{4, 5, 3}) {
		t.Fatalf("Flatten cached outer collection: %v", got)
	}
	calls := 0
	it := VecOf(1, 2, 3).Iter().Map(func(v int) Vec[int] {
		calls++
		return VecOf(v)
	}).Flatten[Vec[int]]()
	if calls != 3 {
		t.Fatalf("measurement calls = %d, want 3", calls)
	}
	if got := it.Collect(Vec_[int]); !slices.Equal(got, []int{1, 2, 3}) || calls != 6 {
		t.Fatalf("Flatten = %v, outer callback calls = %d, want 6", got, calls)
	}
}

func checkIterAgainstLoop(t *testing.T, data []byte) {
	t.Helper()
	src := VecInit(len(data), func(i int) int { return int(int8(data[i])) })
	var want []int
	for _, v := range src {
		if v%3 != 0 {
			for _, x := range []int{v, -v} {
				if x >= 0 {
					want = append(want, x*2+1)
				}
			}
		}
	}
	it := src.Iter().Filter(func(v int) bool { return v%3 != 0 }).
		FlatMap(func(v int) Vec[int] { return VecOf(v, -v) }).
		FilterMap(func(v int) Opt[int] { return Opt_(v*2+1, v >= 0) })
	if got := it.Collect(Vec_[int]); !slices.Equal(got, want) {
		t.Fatalf("pipeline(%v) = %v, want %v", data, got, want)
	}
	groups := src.Iter().Map(func(v int) Vec[int] {
		if v%3 == 0 {
			return nil
		}
		return VecOf(v, -v)
	})
	lazy := groups.FlatMap(func(v Vec[int]) Vec[int] { return v }).Collect(Vec_[int])
	eager := groups.FlatMapEager(func(v Vec[int]) Vec[int] { return v }).Collect(Vec_[int])
	flat := groups.Flatten[Vec[int]]().Collect(Vec_[int])
	if !slices.Equal(lazy, eager) || !slices.Equal(lazy, flat) {
		t.Fatalf("flatten variants differ: lazy=%v eager=%v flatten=%v", lazy, eager, flat)
	}
	wantSum := 17
	for _, v := range want {
		wantSum += v
	}
	if got := it.Reduce(17, func(acc, v int) int { return acc + v }); got != wantSum {
		t.Fatalf("Reduce = %d, want %d", got, wantSum)
	}
}

func TestIterDifferential(t *testing.T) {
	for _, n := range []int{0, 1, 7, 8, 9, 31, 32, 33, 255, 1024} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			data := make([]byte, n)
			for i := range data {
				data[i] = byte(i * 37)
			}
			checkIterAgainstLoop(t, data)
		})
	}
}

func FuzzIterPipeline(f *testing.F) {
	for _, seed := range [][]byte{nil, {0}, {1, 2, 3, 127, 128, 255}, {3, 6, 9}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 4096 {
			t.Skip()
		}
		checkIterAgainstLoop(t, data)
	})
}
