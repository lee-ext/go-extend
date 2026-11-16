package test

import (
	"errors"
	"strings"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

var _resErr = errors.New("boom")

func TestResOk(t *testing.T) {
	r := ResOk(3)
	if !r.IsOk() || r.IsErr() || r.IsNone() {
		t.Fatalf("ResOk states = %v", r)
	}
	if r.Get() != 3 || r.Get_() != 3 || r.GetOr(9) != 3 || r.GetElse(func() int { return 8 }) != 3 {
		t.Errorf("ResOk getters = %v", r)
	}
	if v, err := r.D(); v != 3 || err != nil {
		t.Errorf("ResOk.D = %v,%v", v, err)
	}
	if err := r.Err(); err == nil || !strings.Contains(err.Error(), "ok") {
		t.Errorf("ResOk.Err = %v", err)
	}
	if o := r.ToOpt(); !o.IsSome() || o.Get() != 3 {
		t.Errorf("ResOk.ToOpt = %v", o)
	}
	if _, err := Res_(3, nil).D(); err != nil {
		t.Errorf("Res_ ok err = %v", err)
	}
}

func TestResErr(t *testing.T) {
	r := ResErr[int](_resErr)
	if r.IsOk() || !r.IsErr() || r.IsNone() {
		t.Fatalf("ResErr states = %v", r)
	}
	if v, err := r.D(); v != 0 || !errors.Is(err, _resErr) {
		t.Errorf("ResErr.D = %v,%v", v, err)
	}
	if err := r.Err(); !errors.Is(err, _resErr) {
		t.Errorf("ResErr.Err = %v", err)
	}
	if r.Get_() != 0 || r.GetOr(7) != 7 || r.GetElse(func() int { return 6 }) != 6 {
		t.Error("bad ResErr fallbacks")
	}
	if o := r.ToOpt(); o.IsSome() {
		t.Errorf("ResErr.ToOpt = %v", o)
	}
	if !panics(func() { r.Get() }) {
		t.Error("ResErr.Get should panic")
	}
	if mapped := r.Map(func(i int) Res[string] { return ResOk("x") }); !errors.Is(mapped.Err(), _resErr) {
		t.Errorf("ResErr.Map = %v", mapped)
	}
	if _, err := Res_(3, _resErr).D(); !errors.Is(err, _resErr) {
		t.Errorf("Res_ err = %v", err)
	}
	if def := ResErr[int](nil); def.Err() == nil {
		t.Error("ResErr(nil) should carry a default error")
	}
}

func TestResNone(t *testing.T) {
	r := ResNone[int]()
	if r.IsOk() || !r.IsNone() || !r.IsErr() {
		t.Fatalf("ResNone states = %v", r)
	}
	if v, err := r.D(); v != 0 || err == nil {
		t.Errorf("ResNone.D = %v,%v", v, err)
	}
	if r.Get_() != 0 || r.GetOr(4) != 4 || r.ToOpt().IsSome() {
		t.Error("bad ResNone fallbacks")
	}
	if !panics(func() { r.Get() }) {
		t.Error("ResNone.Get should panic")
	}
}

func TestResUnitAndOpt(t *testing.T) {
	if u := ResUnit(nil); !u.IsOk() || u.Get() != (Unit{}) {
		t.Errorf("ResUnit ok = %v", u)
	}
	if u := ResUnit(_resErr); !errors.Is(u.Err(), _resErr) {
		t.Errorf("ResUnit err = %v", u)
	}
	if r := ResOpt(5, true, nil); r.Get() != 5 {
		t.Errorf("ResOpt some = %v", r)
	}
	if r := ResOpt(5, false, nil); !r.IsNone() {
		t.Errorf("ResOpt none = %v", r)
	}
	if r := ResOpt(5, true, _resErr); !errors.Is(r.Err(), _resErr) {
		t.Errorf("ResOpt err wins = %v", r)
	}
	if (Unit{}).String() != "unit" {
		t.Errorf("Unit.String = %q", (Unit{}).String())
	}
}

func TestResCall(t *testing.T) {
	hits := 0
	ResOk(1).Call(func(int) { hits++ }).ErrCall(func(error) { hits += 10 })
	if hits != 1 {
		t.Errorf("ok Call hits = %d", hits)
	}
	hits = 0
	ResErr[int](_resErr).Call(func(int) { hits++ }).ErrCall(func(error) { hits += 10 })
	if hits != 10 {
		t.Errorf("err ErrCall hits = %d", hits)
	}
	hits = 0
	ResNone[int]().Call(func(int) { hits++ }).ErrCall(func(error) { hits += 10 })
	if hits != 10 {
		t.Errorf("none ErrCall hits = %d", hits)
	}
}

func TestResMap(t *testing.T) {
	got := ResOk(2).Map(func(i int) Res[string] { return ResOk(strings.Repeat("a", i)) })
	if v, err := got.D(); v != "aa" || err != nil {
		t.Errorf("Map = %v,%v", v, err)
	}
	chained := ResOk(1).
		Map(func(i int) Res[int] { return ResOk(i + 1) }).
		Map(func(i int) Res[int] { return ResErr[int](_resErr) }).
		Map(func(i int) Res[int] { return ResOk(i + 100) })
	if !errors.Is(chained.Err(), _resErr) {
		t.Errorf("chained Map = %v", chained)
	}
}

func TestResTry(t *testing.T) {
	ok := ResTry(func() int { return 11 })
	if ok.Get() != 11 || !ok.IsOk() {
		t.Errorf("ResTry ok = %v", ok)
	}
	e := ResTry(func() int { panic(_resErr) })
	if !errors.Is(e.Err(), _resErr) {
		t.Errorf("ResTry error panic = %v", e)
	}
	s := ResTry(func() int { panic("raw") })
	if !strings.Contains(s.Err().Error(), "raw") {
		t.Errorf("ResTry raw panic = %v", s)
	}
	// Go 1.21+ turns panic(nil) into runtime.PanicNilError, so this is an Err.
	if ResTry(func() int { panic(nil) }).IsErr() != true {
		t.Error("ResTry panic(nil) should be err")
	}
}
