package ext

// Actor actor model entity
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
	go func() {
		defer func() {
			r := recover()
			deferFn(r)
			if r != nil {
				go a.receive(deferFn)
			}
		}()
		for fn := range a.ch {
			fn()
		}
	}()
}

// Launch a function to the actor
// If you need to get the returned result, you can use Promise[T] or chan
//
//	p, f := Promise_[T](0)
//	actor.Launch(func() {
//		t := ...
//		f(t)
//	})
//	return p
func (a Actor) Launch(fn func()) {
	a.ch <- fn
}

func ActorAsync[T any](actor Actor, fn func() T) Promise[T] {
	p, f := Promise_[T](0)
	actor.Launch(func() { f(fn()) })
	return p
}

// Close the actor
func (a Actor) Close() {
	close(a.ch)
}
