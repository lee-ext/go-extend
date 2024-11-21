package ext

import (
	"errors"
	"fmt"
	"unsafe"
)

type _DequePin[E any] struct {
	data []E
	head int
	tail int
}

// Deque 环形队列
type Deque[E any] struct {
	*_DequePin[E]
}

func Deque_[E any](cap_ int) Deque[E] {
	return Deque[E]{&_DequePin[E]{make([]E, 0, cap_), 0, 0}}
}

func change(index, cap_, change int) int {
	index += change
	if index < 0 {
		index += cap_
	} else if index >= cap_ {
		index -= cap_
	}
	return index
}

func (d Deque[E]) grow() {
	data := make([]E, max(d.Cap()*2, 1))
	len_ := d.Len()
	if d.head > d.tail {
		headData := d.data[d.head:]
		copy(data, headData)
		copy(data[len(headData):], d.data[:d.tail+1])
		d.data = data
	} else {
		copy(data, d.data[d.head:d.tail+1])
		d.data = data
	}
	d.head = 0
	d.tail = len_ - 1
}

func (d Deque[E]) unsafeEditLen(len int) {
	(*[3]int)(unsafe.Pointer(&d.data))[1] = len
}

func (d Deque[E]) fillLen() bool {
	if d.Empty() {
		d.unsafeEditLen(d.Cap())
		return false
	}
	return true
}

func (d Deque[E]) clearLen() bool {
	if d.head == d.tail {
		d.unsafeEditLen(0)
		return false
	}
	return true
}

func (d Deque[E]) realIndex(index int) int {
	len_ := d.Len()
	if index < 0 || index >= len_ {
		panic(errors.New("index out of range"))
	}
	index += d.head
	if index >= len_ {
		index -= len_
	}
	return index
}

func (d Deque[E]) Get(index int) E {
	return d.data[d.realIndex(index)]
}

func (d Deque[E]) Set(index int, value E) {
	d.data[d.realIndex(index)] = value
}

func (d Deque[E]) PushFront(value E) {
	if d.Fill() {
		d.grow()
	}
	head := d.head
	if d.fillLen() {
		head = change(head, d.Cap(), -1)
	}
	d.data[head] = value
	d.head = head
}

func (d Deque[E]) PushBack(value E) {
	if d.Fill() {
		d.grow()
	}
	tail := d.tail
	if d.fillLen() {
		tail = change(tail, d.Cap(), 1)
	}
	d.data[tail] = value
	d.tail = tail
}

func (d Deque[E]) PopFront() Opt[E] {
	if d.Empty() {
		return Opt_(*new(E), false)
	}
	value := d.data[d.head]
	if d.clearLen() {
		d.data[d.head] = *new(E)
		d.head = change(d.head, d.Cap(), 1)
	}
	return Opt_(value, true)
}

func (d Deque[E]) PopBack() Opt[E] {
	if d.Empty() {
		return Opt_(*new(E), false)
	}
	value := d.data[d.tail]
	if d.clearLen() {
		d.data[d.tail] = *new(E)
		d.tail = change(d.tail, d.Cap(), -1)
	}
	return Opt_(value, true)
}

func (d Deque[E]) Front() Opt[E] {
	if d.Empty() {
		return Opt_(*new(E), false)
	}
	return Opt_(d.data[d.head], true)
}

func (d Deque[E]) Back() Opt[E] {
	if d.Empty() {
		return Opt_(*new(E), false)
	}
	return Opt_(d.data[d.tail], true)
}

func (d Deque[E]) Foreach(fn func(E)) {
	if d.head > d.tail {
		for _, e := range d.data[d.head:] {
			fn(e)
		}
		for _, e := range d.data[:d.tail+1] {
			fn(e)
		}
	} else {
		for _, e := range d.data[d.head : d.tail+1] {
			fn(e)
		}
	}
}

func (d Deque[E]) ToVec() Vec[E] {
	data := Vec_[E](d.Len())
	if d.head > d.tail {
		data.Appends(d.data[d.head:]...)
		data.Appends(d.data[:d.tail+1]...)
	} else {
		data.Appends(d.data[d.head : d.tail+1]...)
	}
	return data
}

func (d Deque[E]) Len() int {
	if d.Empty() {
		return 0
	}
	len_ := d.tail - d.head + 1
	if len_ <= 0 {
		len_ += len(d.data)
	}
	return len_
}

func (d Deque[E]) Cap() int {
	return cap(d.data)
}

func (d Deque[E]) Fill() bool {
	return d.Len() == d.Cap()
}

func (d Deque[E]) Empty() bool {
	return len(d.data) == 0
}

func (d Deque[E]) AppendSelf(e E) Deque[E] {
	d.PushBack(e)
	return d
}

func (d Deque[E]) String() string {
	return fmt.Sprintf("deque%v", d.ToVec())
}
