package ext

import (
	"errors"
	"sync"
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
	locker sync.Mutex
	res    PromiseRes[T]
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
	p.waiter.Wait()
	return p.res
}

func (p Promise[T]) Cancel() bool {
	ok := false
	if p.res.status == _PromisePending {
		p.locker.Lock()
		if p.res.status == _PromisePending {
			p.res.status = _PromiseCanceled
			ok = true
		}
		p.locker.Unlock()
		p.waiter.Done()
	}
	return ok
}

func (p Promise[T]) TryGet() PromiseRes[T] {
	return p.res
}

func Promise_[T any]() (Promise[T], func(T)) {
	p := Promise[T]{&_PromisePin[T]{}}
	p.waiter.Add(1)
	f := func(t T) {
		if p.res.status == _PromisePending {
			p.locker.Lock()
			if p.res.status == _PromisePending {
				p.res.result = t
				p.res.status = _PromiseCompleted
			}
			p.locker.Unlock()
			p.waiter.Done()
		}
	}
	return p, f
}
