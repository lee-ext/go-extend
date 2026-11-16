package test

import (
	"cmp"
	"encoding/json"
	"slices"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestVecConstruct(t *testing.T) {
	if v := Vec_[int](8); v.Len() != 0 || v.Cap() != 8 {
		t.Errorf("Vec_ len=%d cap=%d", v.Len(), v.Cap())
	}
	if v := VecOf(1, 2, 3); !slices.Equal(v, Vec[int]{1, 2, 3}) {
		t.Errorf("VecOf = %v", v)
	}
	if v := VecInit[int](3); !slices.Equal(v, Vec[int]{0, 0, 0}) {
		t.Errorf("VecInit zero = %v", v)
	}
	if v := VecInit(3, func(i int) int { return i * i }); !slices.Equal(v, Vec[int]{0, 1, 4}) {
		t.Errorf("VecInit fn = %v", v)
	}
}

func TestVecAccess(t *testing.T) {
	v := VecOf(1, 2, 3)
	if v.Len() != 3 || v.Empty() {
		t.Fatalf("bad meta len=%d", v.Len())
	}
	if got := v.Get(1); !got.IsSome() || got.Get() != 2 {
		t.Errorf("Get(1) = %v", got)
	}
	if v.Get(3).IsSome() || v.Get(99).IsSome() {
		t.Error("out of range Get should be None")
	}
	if v.First().Get() != 1 || v.Last().Get() != 3 {
		t.Error("bad First/Last")
	}
	e := Vec[int]{}
	if !e.Empty() || e.First().IsSome() || e.Last().IsSome() || e.Pop().IsSome() {
		t.Error("empty vec should yield None")
	}
	if got := VecOf(1, 2).AppendSelf(3); !slices.Equal(got, Vec[int]{1, 2, 3}) {
		t.Errorf("AppendSelf = %v", got)
	}
}

func TestVecMutate(t *testing.T) {
	v := VecOf(1, 2, 3)
	v.Append(4)
	v.Appends(5, 6)
	if !slices.Equal(v, Vec[int]{1, 2, 3, 4, 5, 6}) {
		t.Fatalf("Append(s) = %v", v)
	}
	if got := v.Pop(); got.Get() != 6 || v.Len() != 5 {
		t.Fatalf("Pop = %v, rest = %v", got, v)
	}
	v.Insert(0, 0)
	if !slices.Equal(v, Vec[int]{0, 1, 2, 3, 4, 5}) {
		t.Fatalf("Insert = %v", v)
	}
	v.Replace(1, 3, 9)
	if !slices.Equal(v, Vec[int]{0, 9, 3, 4, 5}) {
		t.Fatalf("Replace = %v", v)
	}
	v.RemoveAt(0)
	v.RemoveRange(1, 3)
	if !slices.Equal(v, Vec[int]{9, 5}) {
		t.Fatalf("Remove = %v", v)
	}
	v.Repeat(2)
	if !slices.Equal(v, Vec[int]{9, 5, 9, 5}) {
		t.Fatalf("Repeat = %v", v)
	}
	v.Grow(8)
	if v.Cap() < v.Len()+8 {
		t.Errorf("Grow cap=%d", v.Cap())
	}
	v.Clip()
	if v.Cap() != v.Len() {
		t.Errorf("Clip cap=%d len=%d", v.Cap(), v.Len())
	}
	v.Clear()
	if !v.Empty() || v.Cap() == 0 {
		t.Errorf("Clear len=%d cap=%d", v.Len(), v.Cap())
	}
}

func TestVecOrder(t *testing.T) {
	v := VecOf(3, 1, 2)
	v.SortFunc(cmp.Compare[int])
	if !slices.Equal(v, Vec[int]{1, 2, 3}) || !v.IsSortedFunc(cmp.Compare[int]) {
		t.Fatalf("SortFunc = %v", v)
	}
	idx, ok := v.BinarySearchFunc(2, cmp.Compare[int])
	if !ok || idx != 1 {
		t.Errorf("BinarySearchFunc = %d,%v", idx, ok)
	}
	if v.MinFunc(cmp.Compare[int]) != 1 || v.MaxFunc(cmp.Compare[int]) != 3 {
		t.Error("bad Min/Max")
	}
	if v.IndexFunc(func(i int) bool { return i == 3 }) != 2 {
		t.Error("bad IndexFunc")
	}
	if !v.ContainsFunc(func(i int) bool { return i == 2 }) {
		t.Error("bad ContainsFunc")
	}
	dup := VecOf(1, 1, 2, 2, 3)
	dup.CompactFunc(func(a, b int) bool { return a == b })
	if !slices.Equal(dup, Vec[int]{1, 2, 3}) {
		t.Errorf("CompactFunc = %v", dup)
	}
	stable := VecOf(T2_(1, "a"), T2_(1, "b"), T2_(0, "c"))
	stable.SortStableFunc(func(a, b T2[int, string]) int { return cmp.Compare(a.V0, b.V0) })
	if stable[0].V1 != "c" || stable[1].V1 != "a" || stable[2].V1 != "b" {
		t.Errorf("SortStableFunc = %v", stable)
	}
}

func TestVecCloneReverseShuffle(t *testing.T) {
	src := VecOf(1, 2, 3)
	cp := src.Clone()
	cp.Append(4)
	if src.Len() != 3 {
		t.Errorf("Clone should not alias src: %v", src)
	}
	src.Reverse()
	if !slices.Equal(src, Vec[int]{3, 2, 1}) {
		t.Errorf("Reverse = %v", src)
	}
	sh := VecInit(50, func(i int) int { return i })
	sh.Shuffle()
	sorted := sh.Clone()
	sorted.SortFunc(cmp.Compare[int])
	if !slices.Equal(sorted, VecInit(50, func(i int) int { return i })) {
		t.Error("Shuffle must keep a permutation")
	}
	if slices.Equal(sh, VecInit(50, func(i int) int { return i })) {
		t.Error("Shuffle should reorder 50 elements")
	}
}

func TestVecViews(t *testing.T) {
	v := VecOf(1, 2, 3)
	var rev []int
	v.ToReverse().ForEach(func(i int) { rev = append(rev, i) })
	if !slices.Equal(rev, []int{3, 2, 1}) {
		t.Errorf("RevVec.ForEach = %v", rev)
	}
	if got := v.ToReverse().Get(0); got.Get() != 3 {
		t.Errorf("RevVec.Get(0) = %v", got)
	}
	if v.ToReverse().Len() != 3 || v.ToReverse().Empty() {
		t.Error("bad RevVec meta")
	}
	count := 0
	v.ToReverse().ForEachWhile(func(int) bool { count++; return count < 2 })
	if count != 2 {
		t.Errorf("RevVec.ForEachWhile count = %d", count)
	}

	var pairs []KV[int, int]
	v.ToIndexed().ForEach(func(kv KV[int, int]) { pairs = append(pairs, kv) })
	if len(pairs) != 3 || pairs[2].K != 2 || pairs[2].V != 3 {
		t.Errorf("IdxVec.ForEach = %v", pairs)
	}
	stop := 0
	v.ToIndexed().ForEachWhile(func(KV[int, int]) bool { stop++; return false })
	if stop != 1 {
		t.Errorf("IdxVec.ForEachWhile stop = %d", stop)
	}
	if v.ToIndexed().Len() != 3 || v.ToIndexed().Empty() {
		t.Error("bad IdxVec meta")
	}
}

func TestVecJSON(t *testing.T) {
	b, err := json.Marshal(VecOf(1, 2, 3))
	if err != nil || string(b) != "[1,2,3]" {
		t.Errorf("Marshal = %s, %v", b, err)
	}
	b, err = json.Marshal(Vec[int](nil))
	if err != nil || string(b) != "[]" {
		t.Errorf("Marshal nil = %s, %v", b, err)
	}
	var v Vec[int]
	if err := json.Unmarshal([]byte(`[4,5]`), &v); err != nil || !slices.Equal(v, Vec[int]{4, 5}) {
		t.Errorf("Unmarshal = %v, %v", v, err)
	}
}

func TestCmpVec(t *testing.T) {
	v := CmpVecOf(3, 1, 2, 2)
	v.Sort()
	if !slices.Equal(v, CmpVec[int]{1, 2, 2, 3}) || !v.IsSorted() {
		t.Fatalf("Sort = %v", v)
	}
	if v.Index(3) != 3 || !v.Contains(2) || v.Contains(9) {
		t.Error("bad Index/Contains")
	}
	if v.Min() != 1 || v.Max() != 3 {
		t.Error("bad Min/Max")
	}
	if idx, ok := v.BinarySearch(2); !ok || idx < 1 || idx > 2 {
		t.Errorf("BinarySearch = %d,%v", idx, ok)
	}
	v.Compact()
	if !slices.Equal(v, CmpVec[int]{1, 2, 3}) {
		t.Errorf("Compact = %v", v)
	}
	if got := v.Get(1).Get(); got != 2 {
		t.Errorf("Get(1) = %d", got)
	}
	if v.Get(3).IsSome() {
		t.Error("out of range Get should be None")
	}
	if v.Pop().Get() != 3 || v.Len() != 2 {
		t.Errorf("Pop = %v", v)
	}
	var rev []int
	v.ToReverse().ForEach(func(i int) { rev = append(rev, i) })
	if !slices.Equal(rev, []int{2, 1}) {
		t.Errorf("ToReverse = %v", rev)
	}
	var idxs []int
	v.ToIndexed().ForEach(func(kv KV[int, int]) { idxs = append(idxs, kv.K) })
	if !slices.Equal(idxs, []int{0, 1}) {
		t.Errorf("ToIndexed = %v", idxs)
	}
	b, err := json.Marshal(CmpVecOf(1, 2))
	if err != nil || string(b) != "[1,2]" {
		t.Errorf("Marshal = %s, %v", b, err)
	}
	if b, _ := json.Marshal(CmpVec[int](nil)); string(b) != "[]" {
		t.Errorf("Marshal nil = %s", b)
	}
	if c := CmpVecInit(2, func(i int) int { return i + 1 }); !slices.Equal(c, CmpVec[int]{1, 2}) {
		t.Errorf("CmpVecInit = %c", c)
	}
	if c := CmpVec_[int](4); c.Len() != 0 || c.Cap() != 4 {
		t.Errorf("CmpVec_ len=%d cap=%d", c.Len(), c.Cap())
	}
}

func TestVecMulti(t *testing.T) {
	m2 := VecM2_(VecOf(1, 2, 3), VecOf("a", "b"))
	if m2.Len() != 2 || m2.Empty() {
		t.Fatalf("VecM2 len = %d", m2.Len())
	}
	if got := m2.ToVec(); len(got) != 2 || got[1].V0 != 2 || got[1].V1 != "b" {
		t.Errorf("VecM2.ToVec = %v", got)
	}
	var sum int
	m2.ForEach(func(p T2[int, string]) { sum += p.V0 })
	if sum != 3 {
		t.Errorf("VecM2.ForEach sum = %d", sum)
	}
	n := 0
	m2.ForEachWhile(func(T2[int, string]) bool { n++; return false })
	if n != 1 {
		t.Errorf("VecM2.ForEachWhile n = %d", n)
	}
	if got := m2.Iter().Collect(Vec_[T2[int, string]]); len(got) != 2 {
		t.Errorf("VecM2.Iter collect = %v", got)
	}

	m3 := VecM3_(VecOf(1, 2), VecOf("a", "b", "c"), VecOf(1.5, 2.5))
	if m3.Len() != 2 || m3.Empty() {
		t.Fatalf("VecM3 len = %d", m3.Len())
	}
	if got := m3.ToVec(); got[0].V2 != 1.5 || got[1].V1 != "b" {
		t.Errorf("VecM3.ToVec = %v", got)
	}
	total := 0.0
	m3.ForEach(func(p T3[int, string, float64]) { total += p.V2 })
	if total != 4.0 {
		t.Errorf("VecM3.ForEach total = %v", total)
	}
	k := 0
	m3.ForEachWhile(func(T3[int, string, float64]) bool { k++; return k < 1 })
	if k != 1 {
		t.Errorf("VecM3.ForEachWhile k = %d", k)
	}
	if VecM2_(Vec[int]{}, VecOf("a")).Empty() != true {
		t.Error("empty VecM2 should be Empty")
	}
}
