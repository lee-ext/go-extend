package ext

type Iterator[E any] interface {
	Foreach(func(E))
	Len() int
	Empty() bool
}

type FromIterator[E, ES any] interface {
	AppendSelf(E) ES
}

// Map 将Iterator[T]转成Vec[R]
func Map[T, R any, TS Iterator[T]](ts TS, fn func(T) R) Vec[R] {
	rs := Vec_[R](ts.Len())
	ts.Foreach(func(t T) {
		rs.Append(fn(t))
	})
	return rs
}

// IntactTo 同Map函数但支持更多类型
func IntactTo[T any, TS Iterator[T], RS FromIterator[T, RS]](
	ts TS, toFn func(int) RS) RS {
	rs := toFn(ts.Len())
	ts.Foreach(func(t T) {
		rs = rs.AppendSelf(t)
	})
	return rs
}

// MapTo 同Map函数但支持更多类型
func MapTo[T, R any, TS Iterator[T], RS FromIterator[R, RS]](
	ts TS, fn func(T) R, toFn func(int) RS) RS {
	rs := toFn(ts.Len())
	ts.Foreach(func(t T) {
		rs = rs.AppendSelf(fn(t))
	})
	return rs
}

// Flatten 将Iterator[Iterator[T]]平铺成Vec[E]
func Flatten[T any, TS Iterator[T], TG Iterator[TS]](tg TG) Vec[T] {
	rs := Vec_[T](Reduce(tg, 0,
		func(l int, r TS) int {
			return l + r.Len()
		}))
	tg.Foreach(func(ts TS) {
		ts.Foreach(rs.Append)
	})
	return rs
}

// FlattenTo 同Flatten函数但支持更多类型
func FlattenTo[T any, TS Iterator[T], TG Iterator[TS], FlatTS FromIterator[T, FlatTS]](
	tg TG, toFn func(int) FlatTS) FlatTS {
	rs := toFn(Reduce(tg, 0,
		func(l int, r TS) int {
			return l + r.Len()
		}))
	tg.Foreach(func(ts TS) {
		ts.Foreach(func(t T) {
			rs = rs.AppendSelf(t)
		})
	})
	return rs
}

// FlatMap 将Iterator[T]平铺成Vec[R]
func FlatMap[T, R any, TS Iterator[T], RS Iterator[R]](ts TS, fn func(T) RS) Vec[R] {
	return Flatten[R](Map(ts, fn))
}

// FlatMapTo 同 FlatMap 函数但支持更多类型
func FlatMapTo[T, R any, TS Iterator[T], RS Iterator[R], FlatRS FromIterator[R, FlatRS]](
	ts TS, fn func(T) RS, toFn func(int) FlatRS) FlatRS {
	return FlattenTo(Map(ts, fn), toFn)
}

// Filter 过滤Iterator[T]中不需要的元素
func Filter[T any, TS Iterator[T]](ts TS, fn func(T) bool) Vec[T] {
	rs := Vec_[T](filterLen(ts.Len()))
	ts.Foreach(func(t T) {
		if fn(t) {
			rs.Append(t)
		}
	})
	return rs
}

// FilterTo 同 Filter 函数但支持更多类型
func FilterTo[T any, TS Iterator[T], RS FromIterator[T, RS]](
	ts TS, fn func(T) bool, toFn func(int) RS) RS {
	rs := toFn(filterLen(ts.Len()))
	ts.Foreach(func(t T) {
		if fn(t) {
			rs = rs.AppendSelf(t)
		}
	})
	return rs
}

// FilterMap 将Iterator[E]转成Iterator[R] 并过滤不需要的元素
func FilterMap[T, R any, TS Iterator[T]](ts TS, fn func(T) Opt[R]) Vec[R] {
	rs := Vec_[R](filterLen(ts.Len()))
	ts.Foreach(func(t T) {
		if r, b := fn(t).D(); b {
			rs.Append(r)
		}
	})
	return rs
}

// FilterMapTo 同FilterMap 函数但支持更多类型
func FilterMapTo[T, R any, TS Iterator[T], RS FromIterator[R, RS]](
	ts TS, fn func(T) Opt[R], toFn func(int) RS) RS {
	rs := toFn(filterLen(ts.Len()))
	ts.Foreach(func(t T) {
		if r, b := fn(t).D(); b {
			rs = rs.AppendSelf(r)
		}
	})
	return rs
}

// Reduce 对Iterator[E]做合并操作 需要一个种子
func Reduce[T, R any, TS Iterator[T]](ts TS, seed R, fn func(R, T) R) R {
	ts.Foreach(func(t T) {
		seed = fn(seed, t)
	})
	return seed
}

// ToDict 分组函数 可以对key映射
func ToDict[K comparable, T any, TS Iterator[T]](ts TS, kFn func(T) K) Dict[K, T] {
	dict := Dict_[K, T](4)
	ts.Foreach(func(t T) {
		dict.Store(kFn(t), t)
	})
	return dict
}

// VToDict 分组函数 可以对key和value映射
func VToDict[K comparable, V, T any, TS Iterator[T]](ts TS, kvFn func(T) (K, V)) Dict[K, V] {
	dict := Dict_[K, V](4)
	ts.Foreach(func(t T) {
		dict.Store(kvFn(t))
	})
	return dict
}

// GroupBy 分组函数 可以对key映射
func GroupBy[K comparable, T any, TS Iterator[T]](ts TS, kFn func(T) K) MDict[K, T] {
	dict := MDict_[K, T](4)
	ts.Foreach(func(t T) {
		dict.Store(kFn(t), t)
	})
	return dict
}

// VGroupBy 分组函数 可以对key和value同时映射
func VGroupBy[K comparable, V, T any, TS Iterator[T]](ts TS, kvFn func(T) (K, V)) MDict[K, V] {
	dict := MDict_[K, V](4)
	ts.Foreach(func(t T) {
		dict.Store(kvFn(t))
	})
	return dict
}

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
