package ext

type Actor struct {
	ch chan func()
}

func Actor_(len_ int, deferFn func(any)) Actor {
	actor := Actor{make(chan func(), len_)}
	go actor.receive(deferFn)
	return actor
}

func (a Actor) receive(deferFn func(any)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				deferFn(recover())
				go a.receive(deferFn)
			}
		}()
		for fn := range a.ch {
			fn()
		}
	}()
}

func (a Actor) Launch(fn func()) {
	a.ch <- fn
}

func (a Actor) Close() {
	close(a.ch)
}
