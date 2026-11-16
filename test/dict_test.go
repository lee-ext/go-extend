package test

import (
	"maps"
	"slices"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestDictConstruct(t *testing.T) {
	if d := Dict_[string, int](8); d.Len() != 0 || !d.Empty() {
		t.Errorf("Dict_ = %v", d)
	}
	d := DictOf(KV_("a", 1), KV_("b", 2))
	if !maps.Equal(d, Dict[string, int]{"a": 1, "b": 2}) {
		t.Errorf("DictOf = %v", d)
	}
	if k, v := KV_("a", 1).D(); k != "a" || v != 1 {
		t.Errorf("KV.D = %v,%v", k, v)
	}
	if s := KV_("a", 1).String(); s != "{a:1}" {
		t.Errorf("KV.String = %q", s)
	}
}

func TestDictStore(t *testing.T) {
	d := Dict_[string, int](0)
	d.Store("a", 1)
	if got := d.Load("a"); !got.IsSome() || got.Get() != 1 {
		t.Fatalf("Load(a) = %v", got)
	}
	if got := d.Load("z"); got.IsSome() || got.Get_() != 0 {
		t.Errorf("Load(z) = %v", got)
	}
	if got := d.LoadOrStore("a", 99); got.Get() != 1 || d.Load("a").Get() != 1 {
		t.Errorf("LoadOrStore existing = %v", got)
	}
	if got := d.LoadOrStore("b", 2); got.IsSome() || d.Load("b").Get() != 2 {
		t.Errorf("LoadOrStore absent = %v", got)
	}
	if got := d.LoadAndDelete("b"); got.Get() != 2 || d.Load("b").IsSome() {
		t.Errorf("LoadAndDelete = %v", got)
	}
	if got := d.LoadAndDelete("zz"); got.IsSome() {
		t.Errorf("LoadAndDelete absent = %v", got)
	}
	d.Delete("a")
	if d.Load("a").IsSome() || !d.Empty() {
		t.Errorf("Delete = %v", d)
	}
	if got := d.AppendSelf(KV_("c", 3)); !maps.Equal(got, d) || d.Load("c").Get() != 3 {
		t.Errorf("AppendSelf = %v", got)
	}
}

func TestDictViews(t *testing.T) {
	d := DictOf(KV_("a", 1), KV_("b", 2))
	keys := d.Keys()
	slices.Sort(keys)
	if !slices.Equal(keys, VecOf("a", "b")) {
		t.Errorf("Keys = %v", keys)
	}
	vals := d.Values()
	slices.Sort(vals)
	if !slices.Equal(vals, VecOf(1, 2)) {
		t.Errorf("Values = %v", vals)
	}
	if v := d.ToVec(); v.Len() != 2 {
		t.Errorf("ToVec = %v", v)
	}
	var sum int
	d.ForEach(func(kv KV[string, int]) { sum += kv.V })
	if sum != 3 {
		t.Errorf("ForEach sum = %d", sum)
	}
	n := 0
	d.ForEachWhile(func(KV[string, int]) bool { n++; return false })
	if n != 1 {
		t.Errorf("ForEachWhile n = %d", n)
	}
	if got := d.Iter().Collect(Dict_[string, int]); !maps.Equal(got, d) {
		t.Errorf("Iter collect = %v", got)
	}
	d.Clear()
	if !d.Empty() {
		t.Errorf("Clear = %v", d)
	}
}

func TestMDictStore(t *testing.T) {
	m := MDict_[string, int](0)
	m.Store("a", 1)
	m.Store("a", 2)
	if got := m.Load("a").Get(); !slices.Equal(got, VecOf(1, 2)) {
		t.Fatalf("Load(a) = %v", got)
	}
	m.MStore("b", 3, 4)
	if got := m.Load("b").Get(); !slices.Equal(got, VecOf(3, 4)) {
		t.Fatalf("Load(b) = %v", got)
	}
	if m.Load("z").IsSome() {
		t.Error("absent Load should be None")
	}
	if old := m.LoadOrMStore("c", 5, 6); old != nil {
		t.Errorf("LoadOrMStore absent old = %v", old)
	}
	if got := m.Load("c").Get(); !slices.Equal(got, VecOf(5, 6)) {
		t.Errorf("LoadOrMStore stored = %v", got)
	}
	if old := m.LoadOrMStore("c", 9); !slices.Equal(old, VecOf(5, 6)) {
		t.Errorf("LoadOrMStore existing old = %v", old)
	}
	if got := m.Load("c").Get(); !slices.Equal(got, VecOf(5, 6)) {
		t.Errorf("LoadOrMStore should not overwrite, got %v", got)
	}
	if got := m.LoadAndDelete("a"); !slices.Equal(got, VecOf(1, 2)) || m.Load("a").IsSome() {
		t.Errorf("LoadAndDelete = %v", got)
	}
	m.Delete("b")
	if m.Load("b").IsSome() {
		t.Error("Delete failed")
	}
	if got := m.AppendSelf(KV_("d", 7)); !slices.Equal(got.Load("d").Get(), VecOf(7)) {
		t.Errorf("AppendSelf = %v", got)
	}
}

func TestMDictViews(t *testing.T) {
	m := MDictOf(KV_("a", VecOf(1, 2)), KV_("b", VecOf(3)))
	if m.Len() != 2 || m.Empty() {
		t.Fatalf("MDictOf len = %d", m.Len())
	}
	keys := m.Keys()
	slices.Sort(keys)
	if !slices.Equal(keys, VecOf("a", "b")) {
		t.Errorf("Keys = %v", keys)
	}
	if v := m.Values(); v.Len() != 2 {
		t.Errorf("Values = %v", v)
	}
	if v := m.ToVec(); v.Len() != 2 {
		t.Errorf("ToVec = %v", v)
	}
	total := 0
	m.ForEach(func(kv KV[string, Vec[int]]) { total += kv.V.Len() })
	if total != 3 {
		t.Errorf("ForEach total = %d", total)
	}
	n := 0
	m.ForEachWhile(func(KV[string, Vec[int]]) bool { n++; return false })
	if n != 1 {
		t.Errorf("ForEachWhile n = %d", n)
	}
	if got := m.Iter().Collect(Vec_[KV[string, Vec[int]]]); got.Len() != 2 {
		t.Errorf("Iter collect = %v", got)
	}
	m.Clear()
	if !m.Empty() || m.Len() != 0 {
		t.Errorf("Clear = %v", m)
	}
}

func TestGroupByRoundTrip(t *testing.T) {
	type item struct {
		Kind string
		No   int
	}
	src := VecOf(item{"a", 1}, item{"b", 2}, item{"a", 3})
	g := GroupBy(src, func(i item) string { return i.Kind })
	if got := g.Load("a").Get(); !slices.Equal(got, VecOf(item{"a", 1}, item{"a", 3})) {
		t.Errorf("GroupBy = %v", g)
	}
	vg := VGroupBy(src, func(i item) (string, int) { return i.Kind, i.No })
	if !slices.Equal(vg.Load("a").Get(), VecOf(1, 3)) {
		t.Errorf("VGroupBy = %v", vg)
	}
}
