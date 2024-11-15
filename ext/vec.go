package ext

import (
	"encoding/json"
	"math/rand/v2"
	"slices"
)

type Vec[E any] []E

func Vec_[E any](cap int) Vec[E] {
	return make([]E, 0, cap)
}

func VecInit[E any](len_ int, fn_ ...func(int) E) Vec[E] {
	vec := make(Vec[E], len_)
	if len(fn_) > 0 {
		fn := fn_[0]
		for i := 0; i < len_; i++ {
			vec[i] = fn(i)
		}
	}
	return vec
}

func VecOf[E any](es ...E) Vec[E] {
	return es
}

func (v Vec[E]) Foreach(fn func(E)) {
	for _, e := range v {
		fn(e)
	}
}

func (v Vec[E]) Len() int {
	return len(v)
}

func (v Vec[E]) Empty() bool {
	return len(v) == 0
}

func (v Vec[E]) Cap() int {
	return cap(v)
}

func (v Vec[E]) Get(index int) Opt[E] {
	if index < v.Len() {
		return Opt[E]{v[index], true}
	}
	return Opt[E]{}
}

func (v Vec[E]) First() Opt[E] {
	if !v.Empty() {
		return Opt[E]{v[0], true}
	}
	return Opt[E]{}
}

func (v Vec[E]) Last() Opt[E] {
	if !v.Empty() {
		return Opt[E]{v[v.Len()-1], true}
	}
	return Opt[E]{}
}

func (v Vec[E]) Reverse() {
	slices.Reverse(v)
}

func (v Vec[E]) Clear() {
	clear(v)
}

func (v *Vec[E]) Append(element E) {
	*v = append(*v, element)
}

func (v *Vec[E]) Appends(elements ...E) {
	*v = append(*v, elements...)
}

func (v *Vec[E]) Insert(index int, elements ...E) {
	*v = slices.Insert(*v, index, elements...)
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

func (v Vec[E]) _AppendSelf(element E) Vec[E] {
	return append(v, element)
}

func (v Vec[E]) IndexForeach(fn func(T2[int, E])) {
	for i, e := range v {
		fn(T2_(i, e))
	}
}

func (v Vec[E]) Shuffle() {
	len_ := v.Len()
	for i := range len_ {
		j := rand.IntN(len_)
		v[i], v[j] = v[j], v[i]
	}
}

func (v Vec[E]) ToRev() RevVec[E] {
	return RevVec[E]{v}
}

type RevVec[E any] struct {
	Vec[E]
}

func (v RevVec[E]) Foreach(fn func(E)) {
	for i := v.Len(); i >= 0; {
		i -= 1
		fn(v.Vec[i])
	}
}

func (v RevVec[E]) Get(index int) Opt[E] {
	return v.Vec.Get(v.Len() - index - 1)
}

var _EmptyJson = []byte("[]")

// MarshalJSON returns m as the JSON encoding of m.
func (v Vec[E]) MarshalJSON() ([]byte, error) {
	if v == nil {
		return _EmptyJson, nil
	}
	return json.Marshal([]E(v))
}
