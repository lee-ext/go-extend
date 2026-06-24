package ext

import (
	"sync"
	"sync/atomic"
)

const (
	_PromisePending    = 0
	_PromiseCompleting = 1
	_PromiseCompleted  = 2
	_PromiseCanceled   = -1
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

func (p Promise[T]) Completed() bool {
	s := p.status.Load()
	return s == _PromiseCompleted || s == _PromiseCompleting
}

func (p Promise[T]) Done() bool {
	s := p.status.Load()
	return s == _PromiseCompleted || s == _PromiseCanceled
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

//func (p Promise[T]) Cancel() bool {
//	if p.status.Load() == _PromisePending {
//		p.locker.Lock()
//		defer p.locker.Unlock()
//		if p.status.Load() == _PromisePending {
//			defer p.waiter.Done()
//			p.status.Store(_PromiseCanceled)
//			return true
//		}
//	}
//	return false
//}

//func (p Promise[T]) Complete(t T) bool {
//	if p.status.Load() == _PromisePending {
//		p.locker.Lock()
//		defer p.locker.Unlock()
//		if p.status.Load() == _PromisePending {
//			defer p.waiter.Done()
//			p.result = t
//			p.status.Store(_PromiseCompleted)
//			return true
//		}
//	}
//	return false
//}
