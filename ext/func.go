package ext

// Map Convert Iterator[T] to Vec[R]
func Map[T, R any, TS Iterator[T]](ts TS, fn func(T) R) Vec[R] {
	rs := Vec_[R](ts.Len())
	ts.ForEach(func(t T) {
		rs.Append(fn(t))
	})
	return rs
}

// IntactTo Convert and keep the internal elements intact
func IntactTo[T any, TS Iterator[T], RS FromIterator[T, RS]](
	ts TS, toFn func(int) RS) RS {
	rs := toFn(ts.Len())
	ts.ForEach(func(t T) {
		rs = rs.AppendSelf(t)
	})
	return rs
}

// MapTo Same as Map but supports more containers
func MapTo[T, R any, TS Iterator[T], RS FromIterator[R, RS]](
	ts TS, fn func(T) R, toFn func(int) RS) RS {
	rs := toFn(ts.Len())
	ts.ForEach(func(t T) {
		rs = rs.AppendSelf(fn(t))
	})
	return rs
}

// MapWhile Convert Iterator[T] to Vec[R] while fn returns false
func MapWhile[T, R any, TS Iterator[T]](
	ts TS, fn func(T) Opt[R]) Vec[R] {
	rs := Vec_[R](filterLen(ts.Len()))
	ts.ForEachWhile(func(t T) bool {
		r, b := fn(t).D()
		if b {
			rs = rs.AppendSelf(r)
		}
		return b
	})
	return rs
}

// MapWhileTo Same as MapWhile but supports more containers
func MapWhileTo[T, R any, TS Iterator[T], RS FromIterator[R, RS]](
	ts TS, fn func(T) Opt[R], toFn func(int) RS) RS {
	rs := toFn(filterLen(ts.Len()))
	ts.ForEachWhile(func(t T) bool {
		r, b := fn(t).D()
		if b {
			rs = rs.AppendSelf(r)
		}
		return b
	})
	return rs
}

// Flatten Flattening the Iterator[T] to Vec[T]
func Flatten[T any, TS Iterator[T], TG Iterator[TS]](tg TG) Vec[T] {
	rs := Vec_[T](Reduce(tg, 0,
		func(l int, r TS) int {
			return l + r.Len()
		}))
	tg.ForEach(func(ts TS) {
		ts.ForEach(rs.Append)
	})
	return rs
}

// FlattenTo Same as Flatten but supports more containers
func FlattenTo[T any, TS Iterator[T], TG Iterator[TS], FlatTS FromIterator[T, FlatTS]](
	tg TG, toFn func(int) FlatTS) FlatTS {
	rs := toFn(Reduce(tg, 0,
		func(l int, r TS) int {
			return l + r.Len()
		}))
	tg.ForEach(func(ts TS) {
		ts.ForEach(func(t T) {
			rs = rs.AppendSelf(t)
		})
	})
	return rs
}

// FlatMap Flattening Iterator[T] to Vec[R]
func FlatMap[T, R any, TS Iterator[T], RS Iterator[R]](ts TS, fn func(T) RS) Vec[R] {
	return Flatten[R](Map(ts, fn))
}

// FlatMapTo Same as FlatMap but supports more containers
func FlatMapTo[T, R any, TS Iterator[T], RS Iterator[R], FlatRS FromIterator[R, FlatRS]](
	ts TS, fn func(T) RS, toFn func(int) FlatRS) FlatRS {
	return FlattenTo(Map(ts, fn), toFn)
}

// Filter Filtering Iterator[T] to Vec[R]
func Filter[T any, TS Iterator[T]](ts TS, fn func(T) bool) Vec[T] {
	rs := Vec_[T](filterLen(ts.Len()))
	ts.ForEach(func(t T) {
		if fn(t) {
			rs.Append(t)
		}
	})
	return rs
}

// FilterTo Same as Filter but supports more containers
func FilterTo[T any, TS Iterator[T], RS FromIterator[T, RS]](
	ts TS, fn func(T) bool, toFn func(int) RS) RS {
	rs := toFn(filterLen(ts.Len()))
	ts.ForEach(func(t T) {
		if fn(t) {
			rs = rs.AppendSelf(t)
		}
	})
	return rs
}

// FilterMap Convert Iterator[T] to Vec[R] and filtering the element
func FilterMap[T, R any, TS Iterator[T]](ts TS, fn func(T) Opt[R]) Vec[R] {
	rs := Vec_[R](filterLen(ts.Len()))
	ts.ForEach(func(t T) {
		if r, b := fn(t).D(); b {
			rs.Append(r)
		}
	})
	return rs
}

// FilterMapTo Same as FilterMap but supports more containers
func FilterMapTo[T, R any, TS Iterator[T], RS FromIterator[R, RS]](
	ts TS, fn func(T) Opt[R], toFn func(int) RS) RS {
	rs := toFn(filterLen(ts.Len()))
	ts.ForEach(func(t T) {
		if r, b := fn(t).D(); b {
			rs = rs.AppendSelf(r)
		}
	})
	return rs
}

// Reduce Iterator[E] requires a seed
func Reduce[T, R any, TS Iterator[T]](ts TS, seed R, fn func(R, T) R) R {
	ts.ForEach(func(t T) {
		seed = fn(seed, t)
	})
	return seed
}

// ToDict Convert to a dictionary and map keys
func ToDict[K comparable, T any, TS Iterator[T]](ts TS, kFn func(T) K) Dict[K, T] {
	return MapTo(ts, func(t T) KV[K, T] {
		return KV_(kFn(t), t)
	}, Dict_[K, T])
}

// VToDict Convert to a dictionary and map keys and values
func VToDict[K comparable, V, T any, TS Iterator[T]](ts TS, kvFn func(T) (K, V)) Dict[K, V] {
	return MapTo(ts, func(t T) KV[K, V] {
		return KV_(kvFn(t))
	}, Dict_[K, V])
}

// GroupBy Grouping function and map keys
func GroupBy[K comparable, T any, TS Iterator[T]](ts TS, kFn func(T) K) MDict[K, T] {
	dict := MDict_[K, T](4)
	ts.ForEach(func(t T) {
		dict.Store(kFn(t), t)
	})
	return dict
}

// VGroupBy Grouping function and map keys and values
func VGroupBy[K comparable, V, T any, TS Iterator[T]](ts TS, kvFn func(T) (K, V)) MDict[K, V] {
	dict := MDict_[K, V](4)
	ts.ForEach(func(t T) {
		dict.Store(kvFn(t))
	})
	return dict
}

// Estimate the length
func filterLen(len_ int) int {
	switch {
	case len_ < 8:
		return len_
	case len_ < 32:
		len2 := len_ / 2
		if len2*2 < len_ {
			len2 += 1
		}
		return len2
	default:
		len4 := len_ / 4
		if len4*4 < len_ {
			len4 += 1
		}
		return len4
	}
}
