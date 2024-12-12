package ext

import "fmt"

type KV[K comparable, V any] struct {
	K K `json:"k"`
	V V `json:"v"`
}

func KV_[K comparable, V any](k K, v V) KV[K, V] {
	return KV[K, V]{k, v}
}

func (kv KV[K, V]) D() (K, V) {
	return kv.K, kv.V
}

func (kv KV[K, V]) String() string {
	return fmt.Sprintf("{%v:%v}", kv.K, kv.V)
}

// Dict Define a generic dictionary
type Dict[K comparable, V any] map[K]V

// Dict_ Create a Dict[K, V] with a specified capacity
func Dict_[K comparable, V any](cap int) Dict[K, V] {
	return make(map[K]V, cap)
}

// DictOf Create a Dict[K, V] that contains the specified keys and values
func DictOf[K comparable, V any](kvs ...KV[K, V]) Dict[K, V] {
	m := make(map[K]V, len(kvs))
	for _, kv := range kvs {
		m[kv.K] = kv.V
	}
	return m
}

// ForEach Traverse the Dict[K, V]
func (d Dict[K, V]) ForEach(fn func(KV[K, V])) {
	for k, v := range d {
		fn(KV_(k, v))
	}
}

// Len Get the number of elements in the Dict[K, V]
func (d Dict[K, V]) Len() int {
	return len(d)
}

// Empty Determine if the Dict[K, V] is empty
func (d Dict[K, V]) Empty() bool {
	return d.Len() == 0
}

// Load Use the key to obtain the value
func (d Dict[K, V]) Load(key K) Opt[V] {
	v, b := d[key]
	return Opt_(v, b)
}

// Store Add key value pairs
func (d Dict[K, V]) Store(key K, value V) {
	d[key] = value
}

// LoadOrStore Add key value pairs, and returns the old value
func (d Dict[K, V]) LoadOrStore(key K, value V) Opt[V] {
	v, b := d[key]
	if !b {
		d[key] = value
	}
	return Opt_(v, b)
}

// Delete Deletes key value pair
func (d Dict[K, V]) Delete(key K) {
	delete(d, key)
}

// LoadAndDelete Deletes key value pair, and returns the old value
func (d Dict[K, V]) LoadAndDelete(key K) Opt[V] {
	v, b := d[key]
	if b {
		delete(d, key)
	}
	return Opt_(v, b)
}

// ToVec Convert a Dict[K, V] to a Vec[KV[K, V]]
func (d Dict[K, V]) ToVec() Vec[KV[K, V]] {
	vec := Vec_[KV[K, V]](d.Len())
	for k, v := range d {
		vec.Append(KV_(k, v))
	}
	return vec
}

// Keys Get all the keys
func (d Dict[K, V]) Keys() Vec[K] {
	vec := Vec_[K](d.Len())
	for k := range d {
		vec.Append(k)
	}
	return vec
}

// Values Get all the Values
func (d Dict[K, V]) Values() Vec[V] {
	vec := Vec_[V](d.Len())
	for _, v := range d {
		vec.Append(v)
	}
	return vec
}

// Clear Empty the Dict[K, V]
func (d Dict[K, V]) Clear() {
	clear(d)
}

// AppendSelf Inserts an element into the Dict[K, V] and returns self
func (d Dict[K, V]) AppendSelf(kv KV[K, V]) Dict[K, V] {
	d[kv.K] = kv.V
	return d
}
