package test

import (
	"encoding/json"
	"maps"
	"slices"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestSetConstruct(t *testing.T) {
	if s := Set_[int](8); s.Len() != 0 || !s.Empty() {
		t.Errorf("Set_ len=%d", s.Len())
	}
	s := SetOf(1, 2, 2, 3)
	if s.Len() != 3 {
		t.Errorf("SetOf should dedup, len=%d", s.Len())
	}
	if !maps.Equal(s, Set[int]{1: {}, 2: {}, 3: {}}) {
		t.Errorf("SetOf = %v", s)
	}
}

func TestSetBasicOps(t *testing.T) {
	s := Set_[int](0)
	s.Insert(1)
	s.Insert(1)
	s.Insert(2)
	if s.Len() != 2 || !s.Contains(1) || s.Contains(3) {
		t.Fatalf("Insert = %v", s)
	}
	s.Remove(1)
	if s.Contains(1) || s.Len() != 1 {
		t.Fatalf("Remove = %v", s)
	}
	if got := s.AppendSelf(7); !got.Contains(7) || got.Len() != 2 {
		t.Errorf("AppendSelf = %v", got)
	}
	var seen []int
	s.ForEach(func(e int) { seen = append(seen, e) })
	slices.Sort(seen)
	if !slices.Equal(seen, []int{2, 7}) {
		t.Errorf("ForEach = %v", seen)
	}
	n := 0
	s.ForEachWhile(func(int) bool { n++; return false })
	if n != 1 {
		t.Errorf("ForEachWhile n = %d", n)
	}
	if v := s.ToVec(); v.Len() != 2 {
		t.Errorf("ToVec = %v", v)
	}
	if it := s.Iter().Collect(Set_[int]); it.Len() != 2 {
		t.Errorf("Iter collect = %v", it)
	}
	s.Clear()
	if !s.Empty() || s.Len() != 0 {
		t.Errorf("Clear = %v", s)
	}
	if SetOf(1).String() != "set[1]" {
		t.Errorf("String = %s", SetOf(1).String())
	}
}

func TestSetAlgebra(t *testing.T) {
	a, b := SetOf(1, 2, 3), SetOf(3, 4)
	if got := a.Or(b); !maps.Equal(got, Set[int]{1: {}, 2: {}, 3: {}, 4: {}}) {
		t.Errorf("Or = %v", got)
	}
	if got := b.Or(a); !maps.Equal(got, Set[int]{1: {}, 2: {}, 3: {}, 4: {}}) {
		t.Errorf("Or reversed = %v", got)
	}
	if got := a.And(b); !maps.Equal(got, Set[int]{3: {}}) {
		t.Errorf("And = %v", got)
	}
	if got := b.And(a); !maps.Equal(got, Set[int]{3: {}}) {
		t.Errorf("And reversed = %v", got)
	}
	if got := a.Sub(b); !maps.Equal(got, Set[int]{1: {}, 2: {}}) {
		t.Errorf("Sub = %v", got)
	}
	if got := b.Sub(a); !maps.Equal(got, Set[int]{4: {}}) {
		t.Errorf("Sub reversed = %v", got)
	}
	if got := a.Xor(b); !maps.Equal(got, Set[int]{1: {}, 2: {}, 4: {}}) {
		t.Errorf("Xor = %v", got)
	}
	if a.Contains(4) || !a.Contains(1) || b.Len() != 2 {
		t.Error("set algebra must not mutate operands")
	}
}

func TestSetAlgebraEdge(t *testing.T) {
	e, a := Set[int]{}, SetOf(1, 2)
	if got := a.Or(e); !maps.Equal(got, a) {
		t.Errorf("Or empty = %v", got)
	}
	if got := a.And(e); got.Len() != 0 {
		t.Errorf("And empty = %v", got)
	}
	if got := a.Sub(e); !maps.Equal(got, a) {
		t.Errorf("Sub empty = %v", got)
	}
	if got := e.Sub(a); got.Len() != 0 {
		t.Errorf("empty Sub = %v", got)
	}
	if got := a.Xor(e); !maps.Equal(got, a) {
		t.Errorf("Xor empty = %v", got)
	}
	if got := a.Xor(maps.Clone(a)); got.Len() != 0 {
		t.Errorf("Xor self = %v", got)
	}
	// force the large-branch of Sub
	big := Set_[int](0)
	for i := range 100 {
		big.Insert(i)
	}
	if got := big.Sub(SetOf(0, 1, 2)); got.Len() != 97 || got.Contains(0) {
		t.Errorf("big Sub len = %d", got.Len())
	}
}

func TestSetJSON(t *testing.T) {
	b, err := json.Marshal(SetOf(5))
	if err != nil || string(b) != "[5]" {
		t.Errorf("Marshal = %s, %v", b, err)
	}
	var s Set[int]
	if err := json.Unmarshal([]byte(`[1,2,2,3]`), &s); err != nil {
		t.Fatalf("Unmarshal err = %v", err)
	}
	if s.Len() != 3 || !s.Contains(2) {
		t.Errorf("Unmarshal = %v", s)
	}
	if err := json.Unmarshal([]byte(`bad`), &s); err == nil {
		t.Error("Unmarshal should fail on bad input")
	}
}
