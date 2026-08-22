package ext

import "fmt"

type Iter[T any] struct {
	ts Iterator[T]
}

func (i Iter[T]) ForEach(fn func(T)) {
	i.ts.ForEach(fn)
}

func (i Iter[T]) ForEachWhile(fn func(T) bool) {
	i.ts.ForEachWhile(fn)
}

func (i Iter[T]) Len() int {
	return i.ts.Len()
}

func (i Iter[T]) Empty() bool {
	return i.ts.Empty()
}

type MapIter[T, R any] struct {
	ts Iterator[T]
	fn func(T) R
}

func (i MapIter[T, R]) ForEach(fn func(R)) {
	i.ts.ForEach(func(t T) {
		fn(i.fn(t))
	})
}

func (i MapIter[T, R]) ForEachWhile(fn func(R) bool) {
	i.ts.ForEachWhile(func(t T) bool {
		return fn(i.fn(t))
	})
}

func (i MapIter[T, R]) Len() int {
	return i.ts.Len()
}

func (i MapIter[T, R]) Empty() bool {
	return i.ts.Empty()
}

type FilterIter[T any] struct {
	ts Iterator[T]
	fn func(T) bool
}

func (i FilterIter[T]) ForEach(fn func(T)) {
	i.ts.ForEach(func(t T) {
		if i.fn(t) {
			fn(t)
		}
	})
}

func (i FilterIter[T]) ForEachWhile(fn func(T) bool) {
	i.ts.ForEachWhile(func(t T) bool {
		if i.fn(t) {
			return fn(t)
		}
		return true
	})
}

func (i FilterIter[T]) Len() int {
	return filterLen(i.ts.Len())
}

func (i FilterIter[T]) Empty() bool {
	return i.ts.Empty()
}

type FilterMapIter[T, R any] struct {
	ts Iterator[T]
	fn func(T) Opt[R]
}

func (i FilterMapIter[T, R]) ForEach(fn func(R)) {
	i.ts.ForEach(func(t T) {
		if r, b := i.fn(t).D(); b {
			fn(r)
		}
	})
}

func (i FilterMapIter[T, R]) ForEachWhile(fn func(R) bool) {
	i.ts.ForEachWhile(func(t T) bool {
		if r, b := i.fn(t).D(); b {
			return fn(r)
		}
		return true
	})
}

func (i FilterMapIter[T, R]) Len() int {
	return filterLen(i.ts.Len())
}

func (i FilterMapIter[T, R]) Empty() bool {
	return i.ts.Empty()
}

type FlattenIter[T any, TS Iterator[T]] struct {
	tg    Iterator[TS]
	total int
}

func (i FlattenIter[T, TS]) ForEach(fn func(T)) {
	i.tg.ForEach(func(ts TS) {
		ts.ForEach(fn)
	})
}

func (i FlattenIter[T, TS]) ForEachWhile(fn func(T) bool) {
	i.tg.ForEachWhile(func(ts TS) bool {
		b := true
		ts.ForEachWhile(func(t T) bool {
			b = fn(t)
			return b
		})
		return b
	})
}

func (i FlattenIter[T, TS]) Len() int {
	return i.total
}

func (i FlattenIter[T, TS]) Empty() bool {
	return i.total == 0
}

type FlatMapIter[T, R any, RS Iterator[R]] struct {
	ts Iterator[T]
	fn func(T) RS
}

func (i FlatMapIter[T, R, RS]) ForEach(fn func(R)) {
	i.ts.ForEach(func(t T) {
		i.fn(t).ForEach(fn)
	})
}

func (i FlatMapIter[T, R, RS]) ForEachWhile(fn func(R) bool) {
	i.ts.ForEachWhile(func(t T) bool {
		b := true
		i.fn(t).ForEachWhile(func(r R) bool {
			b = fn(r)
			return b
		})
		return b
	})
}

func (i FlatMapIter[T, R, RS]) Len() int {
	return i.ts.Len()
}

func (i FlatMapIter[T, R, RS]) Empty() bool {
	return i.ts.Empty()
}

func (i Iter[T]) Map[R any](fn func(T) R) Iter[R] {
	return Iter[R]{ts: MapIter[T, R]{ts: i.ts, fn: fn}}
}

func (i Iter[T]) Filter(fn func(T) bool) Iter[T] {
	return Iter[T]{ts: FilterIter[T]{ts: i.ts, fn: fn}}
}

func (i Iter[T]) FilterMap[R any](fn func(T) Opt[R]) Iter[R] {
	return Iter[R]{ts: FilterMapIter[T, R]{ts: i.ts, fn: fn}}
}

func (i Iter[T]) FlatMap[R any, RS Iterator[R]](fn func(T) RS) Iter[R] {
	return Iter[R]{ts: FlatMapIter[T, R, RS]{ts: i.ts, fn: fn}}
}

func (i Iter[T]) FlatMapEager[R any, RS Iterator[R]](fn func(T) RS) Iter[R] {
	total := 0
	rg := Vec_[RS](i.Len())
	i.ForEach(func(t T) {
		rs := fn(t)
		total += rs.Len()
		rg.Append(rs)
	})
	return Iter[R]{ts: FlattenIter[R, RS]{tg: rg, total: total}}
}

func (i Iter[T]) Flatten[RS Iterator[R], R any]() Iter[R] {
	if rg, b := any(i.ts).(Iterator[RS]); b {
		total := 0
		rg.ForEach(func(rs RS) {
			total += rs.Len()
		})
		return Iter[R]{ts: FlattenIter[R, RS]{tg: rg, total: total}}
	}
	panic(fmt.Errorf("Iter.Flatten: Iter[%T] element is not %T", *new(T), *new(RS)))
}

func (i Iter[T]) Collect[RS FromIterator[T, RS]](toFn func(int) RS) RS {
	rs := toFn(i.Len())
	i.ForEach(func(t T) {
		rs = rs.AppendSelf(t)
	})
	return rs
}

func (i Iter[T]) Reduce[R any](seed R, fn func(R, T) R) R {
	i.ForEach(func(t T) {
		seed = fn(seed, t)
	})
	return seed
}
