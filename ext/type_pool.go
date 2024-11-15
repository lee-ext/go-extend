package ext

import "sync"

type Pool[T any] struct {
	p sync.Pool
}

func Pool_[T any](fn func() T) Pool[T] {
	return Pool[T]{p: sync.Pool{
		New: func() any { return fn() },
	}}
}

func (p *Pool[T]) Put(t T) {
	p.p.Put(t)
}

func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}
