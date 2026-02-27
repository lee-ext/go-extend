package ext

import (
	. "cmp"
	"encoding/json"
	"math/rand/v2"
	"slices"
)

// CmpVec Ordered
type CmpVec[E Ordered] []E

func CmpVec_[E Ordered](cap int) Vec[E] {
	return make(Vec[E], 0, cap)
}

func CmpVecInit[E Ordered](len_ int, fn_ ...func(int) E) Vec[E] {
	vec := make(Vec[E], len_)
	if len(fn_) > 0 {
		fn := fn_[0]
		for i := range len_ {
			vec[i] = fn(i)
		}
	}
	return vec
}

func CmpVecOf[E Ordered](es ...E) Vec[E] {
	return es
}

func (v CmpVec[E]) ForEach(fn func(E)) {
	for _, e := range v {
		fn(e)
	}
}

func (v CmpVec[E]) ForEachWhile(fn func(E) bool) {
	for _, e := range v {
		if !fn(e) {
			return
		}
	}
}

func (v CmpVec[E]) Len() int {
	return len(v)
}

func (v CmpVec[E]) Empty() bool {
	return len(v) == 0
}

func (v CmpVec[E]) AppendSelf(element E) CmpVec[E] {
	return append(v, element)
}

func (v CmpVec[E]) Cap() int {
	return cap(v)
}

func (v CmpVec[E]) Get(index int) Opt[E] {
	if index < v.Len() {
		return Some(v[index])
	}
	return None[E]()
}

func (v CmpVec[E]) First() Opt[E] {
	if v.Empty() {
		return None[E]()
	}
	return Some(v[0])

}

func (v CmpVec[E]) Last() Opt[E] {
	if v.Empty() {
		return None[E]()
	}
	return Some(v[v.Len()-1])
}

func (v *CmpVec[E]) Pop() Opt[E] {
	if v.Empty() {
		return None[E]()
	}
	index := v.Len() - 1
	last := (*v)[index]
	v.RemoveAt(index)
	return Some(last)
}

func (v *CmpVec[E]) Append(element E) {
	*v = append(*v, element)
}

func (v *CmpVec[E]) Appends(elements ...E) {
	*v = append(*v, elements...)
}

func (v CmpVec[E]) Clear() {
	clear(v)
}

func (v *CmpVec[E]) Insert(index int, elements ...E) {
	*v = slices.Insert(*v, index, elements...)
}

func (v *CmpVec[E]) Replace(i, j int, e ...E) {
	*v = slices.Replace(*v, i, j, e...)
}

func (v *CmpVec[E]) Repeat(count int) {
	*v = slices.Repeat(*v, count)
}

func (v *CmpVec[E]) RemoveAt(index int) {
	*v = slices.Delete(*v, index, index+1)
}

func (v *CmpVec[E]) RemoveRange(start, end int) {
	*v = slices.Delete(*v, start, end)
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func (v *CmpVec[E]) Grow(n int) {
	*v = slices.Grow(*v, n)
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func (v *CmpVec[E]) Clip() {
	*v = slices.Clip(*v)
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
// The result may have additional unused capacity.
func (v CmpVec[E]) Clone() CmpVec[E] {
	return slices.Clone(v)
}

func (v CmpVec[E]) Reverse() {
	slices.Reverse(v)
}

func (v CmpVec[E]) Shuffle() {
	len_ := v.Len()
	for i := range len_ {
		j := rand.IntN(len_)
		v[i], v[j] = v[j], v[i]
	}
}

func (v *CmpVec[E]) CompactFunc(eq func(E, E) bool) {
	*v = slices.CompactFunc(*v, eq)
}

func (v CmpVec[E]) IndexFunc(f func(E) bool) int {
	return slices.IndexFunc(v, f)
}

func (v CmpVec[E]) ContainsFunc(f func(E) bool) bool {
	return slices.ContainsFunc(v, f)
}

func (v CmpVec[E]) SortFunc(cmp func(a, b E) int) {
	slices.SortFunc(v, cmp)
}

func (v CmpVec[E]) SortStableFunc(cmp func(a, b E) int) {
	slices.SortStableFunc(v, cmp)
}

func (v CmpVec[E]) IsSortedFunc(cmp func(a, b E) int) {
	slices.IsSortedFunc(v, cmp)
}

func (v CmpVec[E]) BinarySearchFunc(target E, cmp func(a, b E) int) (int, bool) {
	return slices.BinarySearchFunc(v, target, cmp)
}

func (v CmpVec[E]) MaxFunc(cmp func(a, b E) int) E {
	return slices.MaxFunc(v, cmp)
}

func (v CmpVec[E]) MinFunc(cmp func(a, b E) int) E {
	return slices.MinFunc(v, cmp)
}

func (v *CmpVec[E]) Compact() {
	*v = slices.Compact(*v)
}

func (v CmpVec[E]) Index(e E) int {
	return slices.Index(v, e)
}

func (v CmpVec[E]) Contains(e E) bool {
	return slices.Contains(v, e)
}

func (v CmpVec[E]) Sort() {
	slices.Sort(v)
}

func (v CmpVec[E]) IsSorted() {
	slices.IsSorted(v)
}

func (v CmpVec[E]) BinarySearch(target E) (int, bool) {
	return slices.BinarySearch(v, target)
}

func (v CmpVec[E]) Max() E {
	return slices.Max(v)
}

func (v CmpVec[E]) Min() E {
	return slices.Min(v)
}

func (v CmpVec[E]) ToReverse() RevVec[E] {
	return RevVec[E]{Vec[E](v)}
}

func (v CmpVec[E]) ToIndexed() IdxVec[E] {
	return IdxVec[E]{Vec[E](v)}
}

// MarshalJSON returns m as the JSON encoding of m.
func (v CmpVec[E]) MarshalJSON() ([]byte, error) {
	if v == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]E(v))
}
