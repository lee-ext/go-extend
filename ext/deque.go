package ext

import (
	"errors"
	"fmt"
	"iter"
)

// Deque Define a generic ring queue
type Deque[E any] struct {
	*_DequePin[E]
}

type _DequePin[E any] struct {
	data []E
	head int
	tail int
}

// Deque_ Create a Deque[E] with a specified capacity
func Deque_[E any](cap_ int) Deque[E] {
	return Deque[E]{&_DequePin[E]{make([]E, 0, cap_), 0, 0}}
}

// Get Use index to retrieve elements
func (d Deque[E]) Get(index int) E {
	return d.data[d.realIndex(index)]
}

// Set Use index to set elements
func (d Deque[E]) Set(index int, value E) {
	d.data[d.realIndex(index)] = value
}

// PushFront Add an element to the front
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

// PushBack Add an element to the back
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

// PopFront Pop elements from the front
func (d Deque[E]) PopFront() Opt[E] {
	if d.Empty() {
		return None[E]()
	}
	value := d.data[d.head]
	if d.clearLen() {
		d.data[d.head] = *new(E)
		d.head = change(d.head, d.Cap(), 1)
	}
	return Some(value)
}

// PopBack Pop elements from the back
func (d Deque[E]) PopBack() Opt[E] {
	if d.Empty() {
		return None[E]()
	}
	value := d.data[d.tail]
	if d.clearLen() {
		d.data[d.tail] = *new(E)
		d.tail = change(d.tail, d.Cap(), -1)
	}
	return Some(value)
}

// Front View the front element
func (d Deque[E]) Front() Opt[E] {
	if d.Empty() {
		return None[E]()
	}
	return Some(d.data[d.head])
}

// Back View the back element
func (d Deque[E]) Back() Opt[E] {
	if d.Empty() {
		return None[E]()
	}
	return Some(d.data[d.tail])
}

// ForEach Traverse the Deque[E]
func (d Deque[E]) ForEach(fn func(E)) {
	if d.Empty() {
		return
	}
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

func (d Deque[E]) ForEachWhile(fn func(E) bool) {
	if d.Empty() {
		return
	}
	if d.head > d.tail {
		for _, e := range d.data[d.head:] {
			if !fn(e) {
				return
			}
		}
		for _, e := range d.data[:d.tail+1] {
			if !fn(e) {
				return
			}
		}
	} else {
		for _, e := range d.data[d.head : d.tail+1] {
			if !fn(e) {
				return
			}
		}
	}
}

func (d Deque[E]) ToSeq() iter.Seq[E] {
	if d.Empty() {
		return func(func(E) bool) {}
	}
	if d.head > d.tail {
		return func(yield func(E) bool) {
			for _, e := range d.data[d.head:] {
				if !yield(e) {
					return
				}
			}
			for _, e := range d.data[:d.tail+1] {
				if !yield(e) {
					return
				}
			}
		}
	}
	return func(yield func(E) bool) {
		for _, e := range d.data[d.head : d.tail+1] {
			if !yield(e) {
				return
			}
		}
	}
}

func (d Deque[E]) ToSeq2() iter.Seq2[int, E] {
	if d.Empty() {
		return func(func(int, E) bool) {}
	}
	if d.head > d.tail {
		return func(yield func(int, E) bool) {
			for i, e := range d.data[d.head:] {
				if !yield(i, e) {
					return
				}
			}
			offset := len(d.data[d.head:])
			for i, e := range d.data[:d.tail+1] {
				if !yield(offset+i, e) {
					return
				}
			}
		}
	}
	return func(yield func(int, E) bool) {
		for i, e := range d.data[d.head : d.tail+1] {
			if !yield(i, e) {
				return
			}
		}
	}
}

// ToVec Convert a Deque[E] to a Vec[E]
func (d Deque[E]) ToVec() Vec[E] {
	if d.Empty() {
		return Vec[E]{}
	}
	data := Vec_[E](d.Len())
	if d.head > d.tail {
		data.Appends(d.data[d.head:]...)
		data.Appends(d.data[:d.tail+1]...)
	} else {
		data.Appends(d.data[d.head : d.tail+1]...)
	}
	return data
}

// Len Get the number of elements in the Deque[E]
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
	cap_ := d.Cap()
	if cap_ > 0 {
		data := make([]E, cap_*2)
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
	} else {
		d.data = make([]E, 1)
	}
}

/*func (d Deque[E]) unsafeEditLen(len int) {
	(*[3]int)(unsafe.Pointer(&d.data))[1] = len
}*/

func (d Deque[E]) fillLen() bool {
	if d.Empty() {
		//d.unsafeEditLen(d.Cap())
		d.data = d.data[:d.Cap()]
		return false
	}
	return true
}

func (d Deque[E]) clearLen() bool {
	if d.head == d.tail {
		//d.unsafeEditLen(0)
		d.data = d.data[:0]
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
