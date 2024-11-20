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
func (a Actor) Launch(fn func()) {
	a.ch <- fn
}

// Close the actor
func (a Actor) Close() {
	close(a.ch)
}
