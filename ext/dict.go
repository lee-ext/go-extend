package ext

import "fmt"

type KV[K comparable, V any] struct {
	K K `json:"k"`
	V V `json:"v"`
}

func (kv KV[K, V]) String() string {
	return fmt.Sprintf("{%v:%v}", kv.K, kv.V)
}

func KV_[K comparable, V any](k K, v V) KV[K, V] {
	return KV[K, V]{k, v}
}

type Dict[K comparable, V any] map[K]V

func Dict_[K comparable, V any](cap int) Dict[K, V] {
	return make(map[K]V, cap)
}

func (kv KV[K, V]) D() (K, V) {
	return kv.K, kv.V
}

// DictOf 将外部map转为Dict
func DictOf[K comparable, V any](m map[K]V) Dict[K, V] {
	return m
}

func (d Dict[K, V]) Foreach(fn func(KV[K, V])) {
	for k, v := range d {
		fn(KV_(k, v))
	}
}

// Len 求出dict长度
func (d Dict[K, V]) Len() int {
	return len(d)
}

// Empty 判断dict是否为空
func (d Dict[K, V]) Empty() bool {
	return d.Len() == 0
}

// Load 判断key是否存在,并求出对应的值
func (d Dict[K, V]) Load(key K) Opt[V] {
	v, b := d[key]
	return Opt_(v, b)
}

// Store 添加键值对
func (d Dict[K, V]) Store(key K, value V) {
	d[key] = value
}

// LoadOrStore 向Dict中添加键值对，如果key存在，则直接返回
func (d Dict[K, V]) LoadOrStore(key K, value V) Opt[V] {
	v, b := d[key]
	if !b {
		d[key] = value
	}
	return Opt_(v, b)
}

// LoadAndDelete 通过key删除键值对，并且返回v和b，如果key不存在则返回nil，false
func (d Dict[K, V]) LoadAndDelete(key K) Opt[V] {
	v, b := d[key]
	if b {
		delete(d, key)
	}
	return Opt_(v, b)
}

// Delete 删除键值对
func (d Dict[K, V]) Delete(key K) {
	delete(d, key)
}

// ToVec 将dict转为Vec
func (d Dict[K, V]) ToVec() Vec[KV[K, V]] {
	vec := Vec_[KV[K, V]](d.Len())
	for k, v := range d {
		vec.Append(KV_(k, v))
	}
	return vec
}

// Keys 获取所有的key
func (d Dict[K, V]) Keys() Vec[K] {
	vec := Vec_[K](d.Len())
	for k := range d {
		vec.Append(k)
	}
	return vec
}

// Values 获取所有的Values
func (d Dict[K, V]) Values() Vec[V] {
	vec := Vec_[V](d.Len())
	for _, v := range d {
		vec.Append(v)
	}
	return vec
}

func (d Dict[K, V]) Clear() {
	clear(d)
}

func (d Dict[K, V]) _AppendSelf(kv KV[K, V]) Dict[K, V] {
	d[kv.K] = kv.V
	return d
}
