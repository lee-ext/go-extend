package ext

import "time"

// Actor model entity
type Actor struct {
	ch chan func()
}

// Actor_ create a new Actor
func Actor_(cap int, deferFn func(any)) Actor {
	actor := Actor{make(chan func(), cap)}
	go actor.receive(deferFn)
	return actor
}

func (a Actor) receive(deferFn func(any)) {
	defer func() {
		r := recover()
		if r != nil {
			deferFn(r)
			time.AfterFunc(time.Second, func() {
				go a.receive(deferFn)
			})
		}
	}()
	for fn := range a.ch {
		fn()
	}
}

func (a Actor) Launch(fn func()) {
	a.ch <- fn
}

func (a Actor) Await[T any](fn func() T) Promise[T] {
	p := Promise_[T]()
	a.ch <- func() {
		p.Complete(fn())
	}
	return p
}

// Close the actor
func (a Actor) Close() {
	close(a.ch)
}
