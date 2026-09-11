package test

import (
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestTupleDestruct(t *testing.T) {
	if a, b := T2_(1, "x").D(); a != 1 || b != "x" {
		t.Errorf("T2.D = %v,%v", a, b)
	}
	if a, b, c := T3_(1, 2, 3).D(); a+b+c != 6 {
		t.Errorf("T3.D = %v,%v,%v", a, b, c)
	}
	if a, b, c, d := T4_(1, 2, 3, 4).D(); a+b+c+d != 10 {
		t.Errorf("T4.D = %v", a+b+c+d)
	}
	if a, b, c, d, e := T5_(1, 2, 3, 4, 5).D(); a+b+c+d+e != 15 {
		t.Errorf("T5.D = %v", a+b+c+d+e)
	}
	if a, b, c, d, e, f := T6_(1, 2, 3, 4, 5, 6).D(); a+b+c+d+e+f != 21 {
		t.Errorf("T6.D = %v", a+b+c+d+e+f)
	}
	if a, b, c, d, e, f, g := T7_(1, 2, 3, 4, 5, 6, 7).D(); a+b+c+d+e+f+g != 28 {
		t.Errorf("T7.D = %v", a+b+c+d+e+f+g)
	}
	if a, b, c, d, e, f, g, h := T8_(1, 2, 3, 4, 5, 6, 7, 8).D(); a+b+c+d+e+f+g+h != 36 {
		t.Errorf("T8.D = %v", a+b+c+d+e+f+g+h)
	}
	if a, b, c, d, e, f, g, h, i := T9_(1, 2, 3, 4, 5, 6, 7, 8, 9).D(); a+b+c+d+e+f+g+h+i != 45 {
		t.Errorf("T9.D = %v", a+b+c+d+e+f+g+h+i)
	}
}

func TestTupleString(t *testing.T) {
	if s := T2_(1, "x").String(); s != "(1,x)" {
		t.Errorf("T2.String = %q", s)
	}
	if s := T3_(1, 2, 3).String(); s != "(1,2,3)" {
		t.Errorf("T3.String = %q", s)
	}
}

func TestTupleFields(t *testing.T) {
	tup := T3_(1, "a", true)
	if tup.V0 != 1 || tup.V1 != "a" || tup.V2 != true {
		t.Errorf("T3 fields = %+v", tup)
	}
}

func TestKVTuple(t *testing.T) {
	kv := KV_("k", 1)
	if k, v := kv.D(); k != "k" || v != 1 {
		t.Errorf("KV.D = %v,%v", k, v)
	}
	if kv.String() != "{k:1}" {
		t.Errorf("KV.String = %q", kv.String())
	}
	if kv.K != "k" || kv.V != 1 {
		t.Errorf("KV fields = %+v", kv)
	}
}
