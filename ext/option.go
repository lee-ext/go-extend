package ext

import (
	"errors"
	"fmt"
)

const _OptNoneMsg = "option is none"

type Opt[T any] struct {
	v T
	b bool
}

func Opt_[T any](t T, b bool) Opt[T] {
	return Opt[T]{t, b}
}

func PtrToOpt[T any](ptr *T) Opt[T] {
	if ptr != nil {
		return Some(*ptr)
	}
	return None[T]()
}

func Some[T any](t T) Opt[T] {
	return Opt[T]{t, true}
}

func None[T any]() Opt[T] {
	return Opt[T]{}
}

func (o Opt[T]) D() (T, bool) {
	return o.v, o.b
}

func (o Opt[T]) IsSome() bool {
	return o.b
}

func (o Opt[T]) IsNone() bool {
	return !o.b
}

func (o Opt[T]) ToVec() Vec[T] {
	if o.IsSome() {
		return VecOf(o.v)
	}
	return Vec[T]{}
}

// Get if it is `none`, a panic will occur
func (o Opt[T]) Get() T {
	if o.IsSome() {
		return o.v
	}
	panic(errors.New(_OptNoneMsg))
}

// Get_ if it is `none`, return default value
func (o Opt[T]) Get_() T {
	if o.IsSome() {
		return o.v
	}
	return *new(T)
}

func (o Opt[T]) GetOr(t T) T {
	if o.IsSome() {
		return o.v
	}
	return t
}

func (o Opt[T]) GetElse(fn func() T) T {
	if o.IsSome() {
		return o.v
	}
	return fn()
}

func (o *Opt[T]) Take() Opt[T] {
	o_ := *o
	*o = None[T]()
	return o_
}

func (o Opt[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("some(%v)", o.v)
	}
	return "none"
}

// NzOpt Non zero option
type NzOpt[T comparable] struct {
	v T
}

func NzOpt_[T comparable](t T) NzOpt[T] {
	return NzOpt[T]{t}
}

func (o NzOpt[T]) D() (T, bool) {
	return o.v, o.v != *new(T)
}

func (o NzOpt[T]) IsSome() bool {
	return o.v != *new(T)
}

func (o NzOpt[T]) IsNone() bool {
	return o.v == *new(T)
}

func (o NzOpt[T]) ToVec() Vec[T] {
	if o.IsSome() {
		return VecOf(o.v)
	}
	return Vec[T]{}
}

func (o NzOpt[T]) ToOpt() Opt[T] {
	if o.IsSome() {
		return Some(o.v)
	}
	return None[T]()
}

// Get if it is `none`, a panic will occur
func (o NzOpt[T]) Get() T {
	if o.IsSome() {
		return o.v
	}
	panic(errors.New(_OptNoneMsg))
}

func (o NzOpt[T]) GetOr(t T) T {
	if o.IsSome() {
		return o.v
	}
	return t
}

func (o NzOpt[T]) GetElse(fn func() T) T {
	if o.IsSome() {
		return o.v
	}
	return fn()
}

func (o NzOpt[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("some(%v)", o.v)
	}
	return "none"
}
