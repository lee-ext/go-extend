package test

import (
	"errors"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

var _tryErr = errors.New("try err")

func TestTryNoError(t *testing.T) {
	// Try / TryN return values untouched when err is nil.
	Try(nil)
	if v := Try1(1, nil); v != 1 {
		t.Errorf("Try1 = %v", v)
	}
	if a, b := Try2(1, "x", nil); a != 1 || b != "x" {
		t.Errorf("Try2 = %v,%v", a, b)
	}
	if a, b, c := Try3(1, 2, 3, nil); a+b+c != 6 {
		t.Errorf("Try3 = %v", a+b+c)
	}
	if a, b, c, d := Try4(1, 2, 3, 4, nil); a+b+c+d != 10 {
		t.Errorf("Try4 = %v", a+b+c+d)
	}
	if a, b, c, d, e := Try5(1, 2, 3, 4, 5, nil); a+b+c+d+e != 15 {
		t.Errorf("Try5 = %v", a+b+c+d+e)
	}
	if a, b, c, d, e, f := Try6(1, 2, 3, 4, 5, 6, nil); a+b+c+d+e+f != 21 {
		t.Errorf("Try6 = %v", a+b+c+d+e+f)
	}
	if a, b, c, d, e, f, g := Try7(1, 2, 3, 4, 5, 6, 7, nil); a+b+c+d+e+f+g != 28 {
		t.Errorf("Try7 = %v", a+b+c+d+e+f+g)
	}
	if a, b, c, d, e, f, g, h := Try8(1, 2, 3, 4, 5, 6, 7, 8, nil); a+b+c+d+e+f+g+h != 36 {
		t.Errorf("Try8 = %v", a+b+c+d+e+f+g+h)
	}
	if a, b, c, d, e, f, g, h, i := Try9(1, 2, 3, 4, 5, 6, 7, 8, 9, nil); a+b+c+d+e+f+g+h+i != 45 {
		t.Errorf("Try9 = %v", a+b+c+d+e+f+g+h+i)
	}
}

func TestTryPanicsOnError(t *testing.T) {
	if !panics(func() { Try(_tryErr) }) {
		t.Error("Try should panic on error")
	}
	for i, fn := range []func(){
		func() { Try1(0, _tryErr) },
		func() { Try2(0, 0, _tryErr) },
		func() { Try3(0, 0, 0, _tryErr) },
		func() { Try4(0, 0, 0, 0, _tryErr) },
		func() { Try5(0, 0, 0, 0, 0, _tryErr) },
		func() { Try6(0, 0, 0, 0, 0, 0, _tryErr) },
		func() { Try7(0, 0, 0, 0, 0, 0, 0, _tryErr) },
		func() { Try8(0, 0, 0, 0, 0, 0, 0, 0, _tryErr) },
		func() { Try9(0, 0, 0, 0, 0, 0, 0, 0, 0, _tryErr) },
	} {
		if !panics(fn) {
			t.Errorf("Try%d should panic on error", i+1)
		}
	}
}

func TestTryPanicsCarryError(t *testing.T) {
	// the recovered value must be the original error
	var got error
	func() {
		defer func() {
			if r, ok := recover().(error); ok {
				got = r
			}
		}()
		Try(_tryErr)
	}()
	if !errors.Is(got, _tryErr) {
		t.Errorf("recovered = %v", got)
	}
}
