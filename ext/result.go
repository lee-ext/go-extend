package ext

import (
	"errors"
	"fmt"
)

const (
	_ResOkMsg   = "result is ok"
	_ResNoneMsg = "result is none"
)

type Res[T any] struct {
	v any
}

func (r Res[T]) IsOk() bool {
	_, b := r.v.(T)
	return b
}

func (r Res[T]) IsErr() bool {
	return !r.IsOk()
}

func (r Res[T]) IsNone() bool {
	return r.v == nil
}

func (r Res[T]) ToOpt() Opt[T] {
	if t, b := r.v.(T); b {
		return Some(t)
	}
	return None[T]()
}

func (r Res[T]) Get() T {
	switch v := r.v.(type) {
	case T:
		return v
	case error:
		panic(v)
	case nil:
		panic(errors.New(_ResNoneMsg))
	default:
		panic(fmt.Errorf("unknown type: %#v", v))
	}
}

func (r Res[T]) Get_() T {
	switch v := r.v.(type) {
	case T:
		return v
	default:
		return *new(T)
	}
}

func (r Res[T]) GetOr(t T) T {
	switch v := r.v.(type) {
	case T:
		return v
	default:
		return t
	}
}

func (r Res[T]) GetElse(fn func() T) T {
	switch v := r.v.(type) {
	case T:
		return v
	default:
		return fn()
	}
}

func (r Res[T]) Err() error {
	switch v := r.v.(type) {
	case T:
		return errors.New(_ResOkMsg)
	case error:
		return v
	default:
		return errors.New(_ResNoneMsg)
	}
}

func (r Res[T]) Map(fn func(T)) Res[T] {
	if v, b := r.v.(T); b {
		fn(v)
	}
	return r
}

func (r Res[T]) ErrMap(fn func(err error)) Res[T] {
	switch v := r.v.(type) {
	case T:
		break
	case error:
		fn(v)
	default:
		fn(errors.New(_ResNoneMsg))
	}
	return r
}

func (r Res[T]) D() (T, error) {
	switch v := r.v.(type) {
	case T:
		return v, nil
	case error:
		return *new(T), v
	default:
		return *new(T), errors.New(_ResNoneMsg)
	}
}

func Res_[T any](t T, e error) Res[T] {
	if e == nil {
		return Res[T]{t}
	}
	return Res[T]{e}
}

func ResOpt[T any](t T, b bool, e error) Res[T] {
	if e == nil && b {
		return Res[T]{t}
	}
	return Res[T]{e}
}

func ResUnit(e error) Res[Unit] {
	if e == nil {
		return Res[Unit]{Unit{}}
	}
	return Res[Unit]{e}
}

func ResOk[T any](t T) Res[T] {
	return Res[T]{t}
}

func ResErr[T any](e error) Res[T] {
	return Res[T]{e}
}

func ResNone[T any]() Res[T] {
	return Res[T]{nil}
}

func ResMap[T, R any](res Res[T], fn func(T) Res[R]) Res[R] {
	if res.IsOk() {
		return fn(res.Get())
	}
	return ResErr[R](res.Err())
}
