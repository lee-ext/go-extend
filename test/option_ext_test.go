package test

import (
	"encoding/json"
	"slices"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestOptSomeNone(t *testing.T) {
	some, none := Some(7), None[int]()
	if !some.IsSome() || some.IsNone() || none.IsSome() || !none.IsNone() {
		t.Fatal("bad IsSome/IsNone")
	}
	if v, b := some.D(); v != 7 || !b {
		t.Errorf("Some.D = %v,%v", v, b)
	}
	if v, b := none.D(); v != 0 || b {
		t.Errorf("None.D = %v,%v", v, b)
	}
	if none.Get_() != 0 || none.GetOr(5) != 5 || none.GetElse(func() int { return 6 }) != 6 {
		t.Error("bad None getters")
	}
	if some.Get_() != 7 || some.GetOr(5) != 7 || some.GetElse(func() int { return 6 }) != 7 {
		t.Error("bad Some getters")
	}
	if !slices.Equal(some.ToVec(), VecOf(7)) || none.ToVec().Len() != 0 {
		t.Error("bad ToVec")
	}
	if some.String() != "some(7)" || none.String() != "none" {
		t.Errorf("String = %q / %q", some.String(), none.String())
	}
	if !panics(func() { none.Get() }) {
		t.Error("None.Get should panic")
	}
	if Opt_(3, true).Get() != 3 || Opt_(3, false).IsSome() {
		t.Error("bad Opt_")
	}
}

func TestOptTake(t *testing.T) {
	o := Some(1)
	if got := o.Take(); got.Get() != 1 {
		t.Errorf("Take = %v", got)
	}
	if o.IsSome() {
		t.Error("Take should leave None")
	}
	if o.Take().IsSome() {
		t.Error("second Take should be None")
	}
}

func TestPtrToOpt(t *testing.T) {
	v := 9
	if got := PtrToOpt(&v); !got.IsSome() || got.Get() != 9 {
		t.Errorf("PtrToOpt(&v) = %v", got)
	}
	if got := PtrToOpt[int](nil); got.IsSome() {
		t.Errorf("PtrToOpt(nil) = %v", got)
	}
}

func TestNzOpt(t *testing.T) {
	some := NzOpt_(3)
	if !some.IsSome() || some.IsNone() || some.Get() != 3 {
		t.Fatalf("NzOpt some = %v", some)
	}
	if v, b := some.D(); v != 3 || !b {
		t.Errorf("NzOpt.D = %v,%v", v, b)
	}
	none := NzOpt_(0)
	if none.IsSome() || !none.IsNone() {
		t.Fatalf("NzOpt none = %v", none)
	}
	if none.GetOr(4) != 4 || none.GetElse(func() int { return 5 }) != 5 {
		t.Error("bad NzOpt fallbacks")
	}
	if !none.ToOpt().IsNone() || none.ToOpt().IsSome() {
		t.Error("bad NzOpt.ToOpt none")
	}
	if got := some.ToOpt(); got.Get() != 3 {
		t.Errorf("NzOpt.ToOpt = %v", got)
	}
	if none.ToVec().Len() != 0 || some.ToVec().Get(0).Get() != 3 {
		t.Error("bad NzOpt.ToVec")
	}
	if some.String() != "some(3)" || none.String() != "none" {
		t.Errorf("NzOpt String = %q / %q", some.String(), none.String())
	}
	if !panics(func() { none.Get() }) {
		t.Error("NzOpt None.Get should panic")
	}
}

func TestOptU(t *testing.T) {
	o := OptU_[uint32](0, true)
	if !o.IsSome() || o.IsNone() || o.Get() != 0 {
		t.Fatalf("OptU some(0) = %v", o)
	}
	o.Set(42)
	if o.Get() != 42 || o.Get_() != 42 || o.GetOr(1) != 42 || o.GetElse(func() uint32 { return 2 }) != 42 {
		t.Errorf("OptU getters = %v", o)
	}
	n := OptU_[uint32](7, false)
	if n.IsSome() || !n.IsNone() || n.Get_() != 0 || n.GetOr(9) != 9 {
		t.Errorf("OptU none = %v", n)
	}
	if n.GetElse(func() uint32 { return 8 }) != 8 {
		t.Error("bad OptU GetElse")
	}
	if o.String() != "some(42)" || n.String() != "none" {
		t.Errorf("OptU String = %q / %q", o.String(), n.String())
	}
	if !panics(func() { n.Get() }) {
		t.Error("OptU none Get should panic")
	}
	b, err := json.Marshal(o)
	if err != nil || string(b) != "42" {
		t.Errorf("Marshal = %s, %v", b, err)
	}
	if b, _ := json.Marshal(n); string(b) != "null" {
		t.Errorf("Marshal none = %s", b)
	}
	var got OptU[uint32]
	if err := json.Unmarshal([]byte(`17`), &got); err != nil || got.Get() != 17 {
		t.Errorf("Unmarshal = %v, %v", got, err)
	}
	// UnmarshalJSON leaves the receiver untouched on "null".
	var nul OptU[uint32]
	nul.Set(1)
	if err := json.Unmarshal([]byte(`null`), &nul); err != nil || !nul.IsSome() {
		t.Errorf("Unmarshal null keeps prior state = %v, %v", nul, err)
	}
	var fresh OptU[uint32]
	if err := json.Unmarshal([]byte(`null`), &fresh); err != nil || fresh.IsSome() {
		t.Errorf("Unmarshal null on zero value = %v, %v", fresh, err)
	}
	if err := json.Unmarshal([]byte(`x`), &got); err == nil {
		t.Error("Unmarshal should fail on bad input")
	}
	if err := (*OptU[uint32])(nil).UnmarshalJSON([]byte(`1`)); err == nil {
		t.Error("nil pointer Unmarshal should error")
	}
}

func TestOptI(t *testing.T) {
	o := OptI_[int16](0, true)
	if !o.IsSome() || o.Get() != 0 {
		t.Fatalf("OptI some(0) = %v", o)
	}
	o.Set(-5)
	if o.Get() != -5 || o.Get_() != -5 || o.GetOr(1) != -5 || o.GetElse(func() int16 { return 2 }) != -5 {
		t.Errorf("OptI negative = %v", o)
	}
	o.Set(5)
	if o.Get() != 5 {
		t.Errorf("OptI positive = %v", o)
	}
	n := OptI_[int16](7, false)
	if n.IsSome() || !n.IsNone() || n.Get_() != 0 || n.GetOr(3) != 3 {
		t.Errorf("OptI none = %v", n)
	}
	if n.GetElse(func() int16 { return 4 }) != 4 || n.String() != "none" {
		t.Errorf("OptI none fallback = %v", n)
	}
	if o.String() != "some(5)" {
		t.Errorf("OptI String = %q", o.String())
	}
	if !panics(func() { n.Get() }) {
		t.Error("OptI none Get should panic")
	}
	if b, _ := json.Marshal(o); string(b) != "5" {
		t.Errorf("Marshal = %s", b)
	}
	if b, _ := json.Marshal(n); string(b) != "null" {
		t.Errorf("Marshal none = %s", b)
	}
	var got OptI[int16]
	if err := json.Unmarshal([]byte(`-12`), &got); err != nil || got.Get() != -12 {
		t.Errorf("Unmarshal = %v, %v", got, err)
	}
	var nul OptI[int16]
	nul.Set(1)
	if err := json.Unmarshal([]byte(`null`), &nul); err != nil || !nul.IsSome() {
		t.Errorf("Unmarshal null keeps prior state = %v, %v", nul, err)
	}
	var fresh OptI[int16]
	if err := json.Unmarshal([]byte(`null`), &fresh); err != nil || fresh.IsSome() {
		t.Errorf("Unmarshal null on zero value = %v, %v", fresh, err)
	}
	if err := json.Unmarshal([]byte(`zz`), &got); err == nil {
		t.Error("Unmarshal should fail on bad input")
	}
	if err := (*OptI[int16])(nil).UnmarshalJSON([]byte(`1`)); err == nil {
		t.Error("nil pointer Unmarshal should error")
	}
}

func TestOptF(t *testing.T) {
	o := OptF_[float64](0, true)
	if !o.IsSome() || o.Get() != 0 {
		t.Fatalf("OptF some(0) = %v", o)
	}
	o.Set(2.5)
	if o.Get() != 2.5 || o.Get_() != 2.5 || o.GetOr(1) != 2.5 {
		t.Errorf("OptF getters = %v", o)
	}
	o.Set(-1.5)
	if o.Get() != -1.5 {
		t.Errorf("OptF negative = %v", o)
	}
	n := OptF_[float64](7, false)
	if n.IsSome() || !n.IsNone() || n.Get_() != 0 || n.GetOr(3.5) != 3.5 {
		t.Errorf("OptF none = %v", n)
	}
	if n.GetElse(func() float64 { return 4.5 }) != 4.5 || n.String() != "none" {
		t.Errorf("OptF none fallback = %v", n)
	}
	if o.String() != "some(-1.5)" {
		t.Errorf("OptF String = %q", o.String())
	}
	if !panics(func() { n.Get() }) {
		t.Error("OptF none Get should panic")
	}
	if got := o.Opt(); got.Get() != -1.5 || !got.IsSome() {
		t.Errorf("OptF.Opt = %v", got)
	}
	f32 := OptF_[float32](1.25, true)
	if f32.Get() != 1.25 || !f32.IsSome() {
		t.Errorf("OptF float32 = %v", f32)
	}
	if zero32 := OptF_[float32](0, true); !zero32.IsSome() {
		t.Error("OptF float32 some(0) should be Some")
	}
	if b, _ := json.Marshal(o); string(b) != "-1.5" {
		t.Errorf("Marshal = %s", b)
	}
	if b, _ := json.Marshal(n); string(b) != "null" {
		t.Errorf("Marshal none = %s", b)
	}
	var got OptF[float64]
	if err := json.Unmarshal([]byte(`0.75`), &got); err != nil || got.Get() != 0.75 {
		t.Errorf("Unmarshal = %v, %v", got, err)
	}
	var nul OptF[float64]
	nul.Set(1)
	if err := json.Unmarshal([]byte(`null`), &nul); err != nil || !nul.IsSome() {
		t.Errorf("Unmarshal null keeps prior state = %v, %v", nul, err)
	}
	var fresh OptF[float64]
	if err := json.Unmarshal([]byte(`null`), &fresh); err != nil || fresh.IsSome() {
		t.Errorf("Unmarshal null on zero value = %v, %v", fresh, err)
	}
	if err := json.Unmarshal([]byte(`nope`), &got); err == nil {
		t.Error("Unmarshal should fail on bad input")
	}
	if err := (*OptF[float64])(nil).UnmarshalJSON([]byte(`1`)); err == nil {
		t.Error("nil pointer Unmarshal should error")
	}
}
