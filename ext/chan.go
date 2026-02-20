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
	return Opt_(e, b)
}

func (c Chan[E]) TryRecv() Opt[E] {
	select {
	case e, b := <-c:
		return Opt_(e, b)
	default:
		return None[E]()
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
	return Opt_(e, b)
}

const (
	_ChanEmpty    = 0
	_Received     = 1
	_Disconnected = -1
)

type RecvRes[T any] struct {
	v T
	s int8
}

func (r RecvRes[T]) IsChanEmpty() bool {
	return r.s == _ChanEmpty
}

func (r RecvRes[T]) IsReceived() bool {
	return r.s == _Received
}

func (r RecvRes[T]) IsDisconnected() bool {
	return r.s == _Disconnected
}

func (r RecvRes[T]) ToOpt() Opt[T] {
	return Opt_(r.v, r.s == _Received)
}

func (c Receiver[E]) TryRecv() RecvRes[E] {
	select {
	case e, b := <-c:
		if b {
			return RecvRes[E]{e, _Received}
		}
		return RecvRes[E]{s: _Disconnected}
	default:
		return RecvRes[E]{s: _ChanEmpty}
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
