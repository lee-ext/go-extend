package ext

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	_PromisePending    = 0
	_PromiseCompleting = 1
	_PromiseCompleted  = 2
	_PromiseCanceled   = -1
	_PromiseTimedOut   = -2
)

type Promise[T any] struct {
	*_PromisePin[T]
}

type _PromisePin[T any] struct {
	waiter sync.WaitGroup
	status atomic.Int32
	result T
}

func (p Promise[T]) Pending() bool {
	return p.status.Load() == _PromisePending
}

func (p Promise[T]) Canceled() bool {
	return p.status.Load() == _PromiseCanceled
}

func (p Promise[T]) TimedOut() bool {
	return p.status.Load() == _PromiseTimedOut
}

func (p Promise[T]) Completed() bool {
	return p.status.Load() == _PromiseCompleted
}

func (p Promise[T]) Done() bool {
	s := p.status.Load()
	return s == _PromiseCompleted || s == _PromiseCanceled || s == _PromiseTimedOut
}

func (p Promise[T]) Timeout(d time.Duration) *time.Timer {
	return time.AfterFunc(d, func() {
		if p.status.CompareAndSwap(_PromisePending, _PromiseTimedOut) {
			p.waiter.Done()
		}
	})
}

func (p Promise[T]) Cancel() bool {
	if p.status.CompareAndSwap(_PromisePending, _PromiseCanceled) {
		p.waiter.Done()
		return true
	}
	return false
}

func (p Promise[T]) Complete(v T) bool {
	if p.status.CompareAndSwap(_PromisePending, _PromiseCompleting) {
		p.result = v
		p.status.Store(_PromiseCompleted)
		p.waiter.Done()
		return true
	}
	return false
}

func (p Promise[T]) CompleteAfter(d time.Duration, fn func() T) *time.Timer {
	if fn == nil {
		panic("promise: nil completion function")
	}
	return time.AfterFunc(d, func() {
		if p.status.CompareAndSwap(
			_PromisePending,
			_PromiseCompleting,
		) {
			p.result = fn()
			p.status.Store(_PromiseCompleted)
			p.waiter.Done()
		}
	})
}

func (p Promise[T]) Await() Opt[T] {
	p.waiter.Wait()
	return p.TryGet()
}

func (p Promise[T]) TryGet() Opt[T] {
	if p.status.Load() == _PromiseCompleted {
		return Some(p.result)
	}
	return None[T]()
}

func Promise_[T any]() Promise[T] {
	p := Promise[T]{&_PromisePin[T]{}}
	p.waiter.Add(1)
	return p
}
