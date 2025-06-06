package ext

import "time"

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

// Launch a function to the actor
/*If you need to get the returned result, you can use Promise[T] or chan
p := Promise_[T]()
actor.Launch(func() {
	t := ...
	p.Complete(t)
})
return p
*/
func (a Actor) Launch(fn func()) {
	a.ch <- fn
}

// Close the actor
func (a Actor) Close() {
	close(a.ch)
}
