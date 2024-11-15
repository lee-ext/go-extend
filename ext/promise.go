package ext

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	PromisePending  = 0
	PromiseComplete = 1
	PromiseTimeout  = -1
)

type Promise[T any] struct {
	*_PromisePin[T]
}

type _PromisePin[T any] struct {
	waiter  sync.WaitGroup
	pending atomic.Int32
	result  T
}

type PromiseRes[T any] struct {
	Status int8
	result T
}

func (p PromiseRes[T]) IsPending() bool {
	return p.Status == PromisePending
}
func (p PromiseRes[T]) IsTimeout() bool {
	return p.Status == PromiseTimeout
}
func (p PromiseRes[T]) IsComplete() bool {
	return p.Status == PromiseComplete
}

func (p PromiseRes[T]) Get() T {
	if p.IsComplete() {
		return p.result
	}
	panic(errors.New("option is none"))
}

func (p PromiseRes[T]) Get_() T {
	if p.IsComplete() {
		return p.result
	}
	return *new(T)
}

func (p PromiseRes[T]) GetOr(t T) T {
	if p.IsComplete() {
		return p.result
	}
	return t
}

func (p PromiseRes[T]) GetElse(fn func() T) T {
	if p.IsComplete() {
		return p.result
	}
	return fn()
}

// Await 0未完成 1完成 -1超时
func (p Promise[T]) Await() PromiseRes[T] {
	if p.pending.Load() == PromisePending {
		p.waiter.Wait()
	}
	return PromiseRes[T]{int8(p.pending.Load()), p.result}
}

// TryGet 0未完成 1完成 -1超时
func (p Promise[T]) TryGet() PromiseRes[T] {
	return PromiseRes[T]{int8(p.pending.Load()), p.result}
}

func Promise_[T any](timeout time.Duration) (Promise[T], func(T)) {
	p := Promise[T]{&_PromisePin[T]{}}
	p.waiter.Add(1)
	f := func(t T) {
		if p.pending.Load() == PromisePending {
			p.pending.Store(PromiseComplete)
			p.result = t
			p.waiter.Done()
		}
	}
	if timeout > 0 {
		time.AfterFunc(timeout, func() {
			if p.pending.Load() == PromisePending {
				p.pending.Store(PromiseTimeout)
				p.waiter.Done()
			}
		})
	}
	return p, f
}
