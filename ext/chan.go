package ext

type Chan[E any] chan E
type Sender[E any] chan<- E

type Receiver[E any] <-chan E

func Chan_[E any](cap int) Chan[E] {
	return make(chan E, cap)
}

func (c Chan[E]) Split() (sx Sender[E], rx Receiver[E]) {
	ch := (chan E)(c)
	return ch, ch
}

func (c Chan[E]) Send(e E) {
	c <- e
}

func (c Chan[E]) TrySend(e E) bool {
	select {
	case c <- e:
		return true
	default:
		return false
	}
}

func (c Chan[E]) Len() int {
	return len(c)
}

func (c Chan[E]) Cap() int {
	return cap(c)
}

func (c Chan[E]) Full() bool {
	return len(c) == cap(c)
}

func (c Chan[E]) Empty() bool {
	return len(c) == 0
}

func (c Chan[E]) Close() {
	close(c)
}

func (c Chan[E]) AppendSelf(element E) Chan[E] {
	c <- element
	return c
}

func (c Chan[E]) Recv() Opt[E] {
	e, b := <-c
	return Opt[E]{e, b}
}

func (c Chan[E]) TryRecv() Opt[E] {
	select {
	case e, b := <-c:
		return Opt[E]{e, b}
	default:
		return Opt[E]{}
	}
}

func (c Chan[E]) ForEach(fn func(E)) {
	for e := range c {
		fn(e)
	}
}

func (c Chan[E]) ForEachWhile(fn func(E) bool) {
	for e := range c {
		if !fn(e) {
			return
		}
	}
}

func (c Sender[E]) Send(e E) {
	c <- e
}

func (c Sender[E]) TrySend(e E) bool {
	select {
	case c <- e:
		return true
	default:
		return false
	}
}

func (c Sender[E]) Len() int {
	return len(c)
}

func (c Sender[E]) Cap() int {
	return cap(c)
}

func (c Sender[E]) Full() bool {
	return len(c) == cap(c)
}

func (c Sender[E]) Empty() bool {
	return len(c) == 0
}

func (c Sender[E]) Close() {
	close(c)
}

func (c Sender[E]) AppendSelf(element E) Sender[E] {
	c <- element
	return c
}

func (c Receiver[E]) Recv() Opt[E] {
	e, b := <-c
	return Opt[E]{e, b}
}

func (c Receiver[E]) TryRecv() Opt[E] {
	select {
	case e, b := <-c:
		return Opt[E]{e, b}
	default:
		return Opt[E]{}
	}
}

func (c Receiver[E]) ForEach(fn func(E)) {
	for e := range c {
		fn(e)
	}
}

func (c Receiver[E]) ForEachWhile(fn func(E) bool) {
	for e := range c {
		if !fn(e) {
			return
		}
	}
}

func (c Receiver[E]) Len() int {
	return len(c)
}

func (c Receiver[E]) Cap() int {
	return cap(c)
}

func (c Receiver[E]) Full() bool {
	return len(c) == cap(c)
}

func (c Receiver[E]) Empty() bool {
	return len(c) == 0
}
