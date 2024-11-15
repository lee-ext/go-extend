package ext

type Foreach[E any] interface {
	Foreach(func(E))
	Len() int
	Empty() bool
}

type _AppendSelf[E, ES any] interface {
	_AppendSelf(E) ES
}

// Map 将Foreach[T]转成Vec[R]
func Map[T, R any, TS Foreach[T]](ts TS, fn func(T) R) Vec[R] {
	rs := Vec_[R](ts.Len())
	ts.Foreach(func(t T) {
		rs.Append(fn(t))
	})
	return rs
}

// IndexMap 将Vec[T]转成Vec[R]
func IndexMap[T, R any](vec Vec[T], fn func(int, T) R) Vec[R] {
	rs := Vec_[R](vec.Len())
	for i, t := range vec {
		rs.Append(fn(i, t))
	}
	return rs
}

// MapTo 同 Map函数但支持更多类型
func MapTo[T, R any, TS Foreach[T], RS _AppendSelf[R, RS]](
	ts TS, fn func(T) R, toFn func(int) RS) RS {
	rs := toFn(ts.Len())
	ts.Foreach(func(t T) {
		rs = rs._AppendSelf(fn(t))
	})
	return rs
}

// Flatten 将Foreach[Foreach[T]]平铺成Vec[E]
func Flatten[T any, TS Foreach[T], TG Foreach[TS]](tg TG) Vec[T] {
	rs := Vec_[T](Reduce(tg, 0,
		func(l int, r TS) int {
			return l + r.Len()
		}))
	tg.Foreach(func(ts TS) {
		ts.Foreach(rs.Append)
	})
	return rs
}

// FlattenTo 同 Flatten函数但支持更多类型
func FlattenTo[T any, TS Foreach[T], TG Foreach[TS], FlatTS _AppendSelf[T, FlatTS]](
	tg TG, toFn func(int) FlatTS) FlatTS {
	rs := toFn(Reduce(tg, 0,
		func(l int, r TS) int {
			return l + r.Len()
		}))
	tg.Foreach(func(ts TS) {
		ts.Foreach(func(t T) {
			rs = rs._AppendSelf(t)
		})
	})
	return rs
}

// FlatMap 将Foreach[T]平铺成Vec[R]
func FlatMap[T, R any, TS Foreach[T], RS Foreach[R]](ts TS, fn func(T) RS) Vec[R] {
	return Flatten[R](Map(ts, fn))
}

// IndexFlatMap 将Vec[T]平铺成Vec[R]
func IndexFlatMap[T, R any, RS Foreach[R]](vec Vec[T], fn func(int, T) RS) Vec[R] {
	rs := Vec_[RS](vec.Len())
	for i, t := range vec {
		rs.Append(fn(i, t))
	}
	return Flatten[R](rs)
}

// FlatMapTo 同 FlatMap 函数但支持更多类型
func FlatMapTo[T, R any, TS Foreach[T], RS Foreach[R], FlatRS _AppendSelf[R, FlatRS]](
	ts TS, fn func(T) RS, toFn func(int) FlatRS) FlatRS {
	return FlattenTo(Map(ts, fn), toFn)
}

// Filter 过滤Foreach[T]中不需要的元素
func Filter[T any, TS Foreach[T]](ts TS, fn func(T) bool) Vec[T] {
	rs := Vec_[T](filterLen(ts.Len()))
	ts.Foreach(func(t T) {
		if fn(t) {
			rs.Append(t)
		}
	})
	return rs
}

// IndexFilter 过滤Vec[E]中不需要的元素
func IndexFilter[T any](vec Vec[T], fn func(int, T) bool) Vec[T] {
	rs := Vec_[T](filterLen(vec.Len()))
	for i, t := range vec {
		if fn(i, t) {
			rs.Append(t)
		}
	}
	return rs
}

// FilterTo 同 Filter 函数但支持更多类型
func FilterTo[T any, TS Foreach[T], RS _AppendSelf[T, RS]](
	ts TS, fn func(T) bool, toFn func(int) RS) RS {
	rs := toFn(filterLen(ts.Len()))
	ts.Foreach(func(t T) {
		if fn(t) {
			rs = rs._AppendSelf(t)
		}
	})
	return rs
}

// FilterMap 将Vec[E]转成Vec[RangeTo] 并过滤不需要的元素
func FilterMap[T, R any, TS Foreach[T]](ts TS, fn func(T) Opt[R]) Vec[R] {
	rs := Vec_[R](filterLen(ts.Len()))
	ts.Foreach(func(t T) {
		if r, b := fn(t).D(); b {
			rs.Append(r)
		}
	})
	return rs
}

// IndexFilterMap 将Vec[E]转成Vec[RangeTo] 并过滤不需要的元素
func IndexFilterMap[T, R any](vec Vec[T], fn func(int, T) Opt[R]) Vec[R] {
	rs := Vec_[R](filterLen(vec.Len()))
	for i, t := range vec {
		if r, b := fn(i, t).D(); b {
			rs.Append(r)
		}
	}
	return rs
}

// FilterMapTo 同 FilterMap 函数但支持更多类型
func FilterMapTo[T, R any, TS Foreach[T], RS _AppendSelf[R, RS]](
	ts TS, fn func(T) Opt[R], toFn func(int) RS) RS {
	rs := toFn(filterLen(ts.Len()))
	ts.Foreach(func(t T) {
		if r, b := fn(t).D(); b {
			rs = rs._AppendSelf(r)
		}
	})
	return rs
}

// Reduce 对Vec[E]做合并操作 需要一个种子
func Reduce[T, R any, TS Foreach[T]](ts TS, seed R, fn func(R, T) R) R {
	ts.Foreach(func(t T) {
		seed = fn(seed, t)
	})
	return seed
}

// ToDict 分组函数 可以对key映射
func ToDict[K comparable, T any, TS Foreach[T]](ts TS, kFn func(T) K) Dict[K, T] {
	dict := Dict_[K, T](4)
	ts.Foreach(func(t T) {
		dict.Store(kFn(t), t)
	})
	return dict
}

// VToDict 分组函数 可以对key和value映射
func VToDict[K comparable, V, T any, TS Foreach[T]](ts TS, kvFn func(T) (K, V)) Dict[K, V] {
	dict := Dict_[K, V](4)
	ts.Foreach(func(t T) {
		dict.Store(kvFn(t))
	})
	return dict
}

// GroupBy 分组函数 可以对key映射
func GroupBy[K comparable, T any, TS Foreach[T]](ts TS, kFn func(T) K) MDict[K, T] {
	dict := MDict_[K, T](4)
	ts.Foreach(func(t T) {
		dict.Store(kFn(t), t)
	})
	return dict
}

// VGroupBy 分组函数 可以对key和value同时映射
func VGroupBy[K comparable, V, T any, TS Foreach[T]](ts TS, kvFn func(T) (K, V)) MDict[K, V] {
	dict := MDict_[K, V](4)
	ts.Foreach(func(t T) {
		dict.Store(kvFn(t))
	})
	return dict
}

// FollowSort 跟随排序
/*func FollowSort[O comparable, E any](orders Vec[O], vec Vec[E], kFn func(E) O) Vec[E] {
	return MapFilter(orders, ToDict(vec, kFn).Load)
}*/

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

func Init_[T any]() T {
	return *new(T)
}
