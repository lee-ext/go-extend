package ext

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	_PromisePending   = 0
	_PromiseAssigning = 1
	_PromiseCompleted = 2
	_PromiseCanceled  = -1
)

type Promise[T any] struct {
	*_PromisePin[T]
}

type _PromisePin[T any] struct {
	waiter sync.WaitGroup
	status atomic.Int32
	result T
}

type PromiseRes[T any] struct {
	status int8
	result T
}

func (p PromiseRes[T]) IsPending() bool {
	return p.status == _PromisePending || p.status == _PromiseAssigning
}
func (p PromiseRes[T]) IsCanceled() bool {
	return p.status == _PromiseCanceled
}
func (p PromiseRes[T]) IsCompleted() bool {
	return p.status == _PromiseCompleted
}

func (p PromiseRes[T]) Get() T {
	if p.IsCompleted() {
		return p.result
	}
	panic(errors.New("option is none"))
}

func (p PromiseRes[T]) Get_() T {
	if p.IsCompleted() {
		return p.result
	}
	return *new(T)
}

func (p PromiseRes[T]) GetOr(t T) T {
	if p.IsCompleted() {
		return p.result
	}
	return t
}

func (p PromiseRes[T]) GetElse(fn func() T) T {
	if p.IsCompleted() {
		return p.result
	}
	return fn()
}

func (p Promise[T]) Await() PromiseRes[T] {
	if p.status.Load() == _PromisePending {
		p.waiter.Wait()
	}
	return PromiseRes[T]{int8(p.status.Load()), p.result}
}

func (p Promise[T]) Cancel() bool {
	if p.status.Load() == _PromisePending {
		p.status.Store(_PromiseCanceled)
		p.waiter.Done()
		return true
	}
	return false
}

func (p Promise[T]) TryGet() PromiseRes[T] {
	return PromiseRes[T]{int8(p.status.Load()), p.result}
}

func Promise_[T any]() (Promise[T], func(T)) {
	p := Promise[T]{&_PromisePin[T]{}}
	p.waiter.Add(1)
	f := func(t T) {
		if p.status.Load() == _PromisePending {
			p.status.Store(_PromiseAssigning)
			p.result = t
			p.status.Store(_PromiseCompleted)
			p.waiter.Done()
		}
	}
	return p, f
}
