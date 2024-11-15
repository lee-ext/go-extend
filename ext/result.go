package ext

import (
	"errors"
)

type CuRes[T any] struct {
	v any
}

func (r CuRes[T]) IsOk() bool {
	_, b := r.v.(T)
	return b
}

func (r CuRes[T]) IsErr() bool {
	return !r.IsOk()
}

func (r CuRes[T]) IsNone() bool {
	return r.v == nil
}

func (r CuRes[T]) ToOpt() Opt[T] {
	switch v := r.v.(type) {
	case T:
		return Some(v)
	default:
		return None[T]()
	}
}

func (r CuRes[T]) Get() T {
	switch v := r.v.(type) {
	case T:
		return v
	case error:
		panic(v)
	default:
		panic(errors.New("data is none"))
	}
}

func (r CuRes[T]) Get_() T {
	switch v := r.v.(type) {
	case T:
		return v
	default:
		return *new(T)
	}
}

func (r CuRes[T]) GetOr(t T) T {
	switch v := r.v.(type) {
	case T:
		return v
	default:
		return t
	}
}

func (r CuRes[T]) GetElse(fn func() T) T {
	switch v := r.v.(type) {
	case T:
		return v
	default:
		return fn()
	}
}

func (r CuRes[T]) GetErr() error {
	switch v := r.v.(type) {
	case T:
		return errors.New("result is ok")
	case error:
		return v
	default:
		return errors.New("data is none")
	}
}

func (r CuRes[T]) Map(fn func(T)) CuRes[T] {
	if v, b := r.v.(T); b {
		fn(v)
	}
	return r
}

func (r CuRes[T]) ErrMap(fn func(err error)) CuRes[T] {
	switch v := r.v.(type) {
	case T:
		break
	case error:
		fn(v)
	}
	return r
}

func (r CuRes[T]) D() (T, error) {
	switch v := r.v.(type) {
	case T:
		return v, nil
	case error:
		return *new(T), v
	default:
		return *new(T), errors.New("data is none")
	}
}

func CuRes_[T any](t T, e error) CuRes[T] {
	if e == nil {
		return CuRes[T]{t}
	}
	return CuRes[T]{e}
}

func CuResOpt[T any](t T, b bool, e error) CuRes[T] {
	if e == nil && b {
		return CuRes[T]{t}
	}
	return CuRes[T]{e}
}

func CuResUnit(e error) CuRes[Unit] {
	if e == nil {
		return CuRes[Unit]{Unit{}}
	}
	return CuRes[Unit]{e}
}

func CuOk[T any](t T) CuRes[T] {
	return CuRes[T]{t}
}

func CuErr[T any](e error) CuRes[T] {
	return CuRes[T]{e}
}

func CuNone[T any]() CuRes[T] {
	return CuRes[T]{errors.New("data is none")}
}

func CuResTo[T, R any](res CuRes[T], fn func(T) R) CuRes[R] {
	if res.IsOk() {
		return CuOk(fn(res.Get()))
	}
	return CuErr[R](res.GetErr())
}
