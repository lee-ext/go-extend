package ext

import (
	"encoding/json"
	"math/rand/v2"
	"slices"
)

type Vec[E any] []E

func Vec_[E any](cap int) Vec[E] {
	return make(Vec[E], 0, cap)
}

func VecInit[E any](len_ int, fn_ ...func(int) E) Vec[E] {
	vec := make(Vec[E], len_)
	if len(fn_) > 0 {
		fn := fn_[0]
		for i := range len_ {
			vec[i] = fn(i)
		}
	}
	return vec
}

func VecOf[E any](es ...E) Vec[E] {
	return es
}

func (v Vec[E]) ForEach(fn func(E)) {
	for _, e := range v {
		fn(e)
	}
}

func (v Vec[E]) ForEachWhile(fn func(E) bool) {
	for _, e := range v {
		if !fn(e) {
			return
		}
	}
}

func (v Vec[E]) Len() int {
	return len(v)
}

func (v Vec[E]) Empty() bool {
	return len(v) == 0
}

func (v Vec[E]) AppendSelf(element E) Vec[E] {
	return append(v, element)
}

func (v Vec[E]) Cap() int {
	return cap(v)
}

func (v Vec[E]) Get(index int) Opt[E] {
	if index < v.Len() {
		return Some(v[index])
	}
	return None[E]()
}

func (v Vec[E]) First() Opt[E] {
	if v.Empty() {
		return None[E]()
	}
	return Some(v[0])

}

func (v Vec[E]) Last() Opt[E] {
	if v.Empty() {
		return None[E]()
	}
	return Some(v[v.Len()-1])
}

func (v *Vec[E]) Pop() Opt[E] {
	if v.Empty() {
		return None[E]()
	}
	index := v.Len() - 1
	last := (*v)[index]
	v.RemoveAt(index)
	return Some(last)
}

func (v *Vec[E]) Append(element E) {
	*v = append(*v, element)
}

func (v *Vec[E]) Appends(elements ...E) {
	*v = append(*v, elements...)
}

func (v Vec[E]) Clear() {
	clear(v)
}

func (v *Vec[E]) Insert(index int, elements ...E) {
	*v = slices.Insert(*v, index, elements...)
}

func (v *Vec[E]) Replace(i, j int, e ...E) {
	*v = slices.Replace(*v, i, j, e...)
}

func (v *Vec[E]) Repeat(count int) {
	*v = slices.Repeat(*v, count)
}

func (v *Vec[E]) RemoveAt(index int) {
	*v = slices.Delete(*v, index, index+1)
}

func (v *Vec[E]) RemoveRange(start, end int) {
	*v = slices.Delete(*v, start, end)
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func (v *Vec[E]) Grow(n int) {
	*v = slices.Grow(*v, n)
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func (v *Vec[E]) Clip() {
	*v = slices.Clip(*v)
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
// The result may have additional unused capacity.
func (v Vec[E]) Clone() Vec[E] {
	return slices.Clone(v)
}

func (v Vec[E]) Reverse() {
	slices.Reverse(v)
}

func (v Vec[E]) Shuffle() {
	len_ := v.Len()
	for i := range len_ {
		j := rand.IntN(len_)
		v[i], v[j] = v[j], v[i]
	}
}

func (v *Vec[E]) CompactFunc(eq func(E, E) bool) {
	*v = slices.CompactFunc(*v, eq)
}

func (v Vec[E]) IndexFunc(f func(E) bool) int {
	return slices.IndexFunc(v, f)
}

func (v Vec[E]) ContainsFunc(f func(E) bool) bool {
	return slices.ContainsFunc(v, f)
}

func (v Vec[E]) SortFunc(cmp func(a, b E) int) {
	slices.SortFunc(v, cmp)
}

func (v Vec[E]) SortStableFunc(cmp func(a, b E) int) {
	slices.SortStableFunc(v, cmp)
}

func (v Vec[E]) IsSortedFunc(cmp func(a, b E) int) {
	slices.IsSortedFunc(v, cmp)
}

func (v Vec[E]) BinarySearchFunc(target E, cmp func(a, b E) int) (int, bool) {
	return slices.BinarySearchFunc(v, target, cmp)
}

func (v Vec[E]) MaxFunc(cmp func(a, b E) int) E {
	return slices.MaxFunc(v, cmp)
}

func (v Vec[E]) MinFunc(cmp func(a, b E) int) E {
	return slices.MinFunc(v, cmp)
}

type RevVec[E any] struct {
	Vec[E]
}

func (v Vec[E]) ToReverse() RevVec[E] {
	return RevVec[E]{v}
}

func (v RevVec[E]) ForEach(fn func(E)) {
	for i := v.Len(); i > 0; {
		i -= 1
		fn(v.Vec[i])
	}
}

func (v RevVec[E]) ForEachWhile(fn func(E) bool) {
	for i := v.Len(); i > 0; {
		i -= 1
		if !fn(v.Vec[i]) {
			break
		}
	}
}

func (v RevVec[E]) Get(index int) Opt[E] {
	return v.Vec.Get(v.Len() - index - 1)
}

type IdxVec[E any] struct {
	Vec[E]
}

func (v Vec[E]) ToIndexed() IdxVec[E] {
	return IdxVec[E]{v}
}

func (v IdxVec[E]) ForEach(fn func(KV[int, E])) {
	for i, e := range v.Vec {
		fn(KV_(i, e))
	}
}

func (v IdxVec[E]) ForEachWhile(fn func(KV[int, E]) bool) {
	for i, e := range v.Vec {
		if !fn(KV_(i, e)) {
			return
		}
	}
}

// MarshalJSON returns m as the JSON encoding of m.
func (v Vec[E]) MarshalJSON() ([]byte, error) {
	if v == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]E(v))
}
