package ext

type MDict[K comparable, V any] map[K]Vec[V]

func MDict_[K comparable, V any](cap int) MDict[K, V] {
	return make(map[K]Vec[V], cap)
}

func MDictOf[K comparable, V any](kvs ...KV[K, Vec[V]]) MDict[K, V] {
	m := make(map[K]Vec[V], len(kvs))
	for _, kv := range kvs {
		m[kv.K] = kv.V
	}
	return m
}

func (d MDict[K, V]) ForEach(fn func(KV[K, Vec[V]])) {
	for k, v := range d {
		fn(KV_(k, v))
	}
}

func (d MDict[K, V]) Len() int {
	return len(d)
}

func (d MDict[K, V]) Empty() bool {
	return d.Len() == 0
}

func (d MDict[K, V]) Load(key K) Opt[Vec[V]] {
	v, b := d[key]
	return Opt_(v, b)
}

func (d MDict[K, V]) Store(key K, value V) {
	d[key] = append(d[key], value)
}

func (d MDict[K, V]) MStore(key K, values ...V) {
	d[key] = append(d[key], values...)
}

func (d MDict[K, V]) LoadOrMStore(key K, values ...V) Vec[V] {
	v, b := d[key]
	if !b {
		d[key] = values
	}
	return v
}

func (d MDict[K, V]) LoadAndDelete(key K) Vec[V] {
	v, b := d[key]
	if b {
		delete(d, key)
	}
	return v
}

func (d MDict[K, V]) Delete(key K) {
	delete(d, key)
}

func (d MDict[K, V]) ToVec() Vec[KV[K, Vec[V]]] {
	vec := Vec_[KV[K, Vec[V]]](d.Len())
	for k, v := range d {
		vec.Append(KV_(k, v))
	}
	return vec
}

func (d MDict[K, V]) Keys() Vec[K] {
	vec := Vec_[K](d.Len())
	for k := range d {
		vec.Append(k)
	}
	return vec
}

func (d MDict[K, V]) Values() Vec[Vec[V]] {
	vec := Vec_[Vec[V]](d.Len())
	for _, v := range d {
		vec.Append(v)
	}
	return vec
}

func (d MDict[K, V]) Clear() {
	clear(d)
}

func (d MDict[K, V]) AppendSelf(kv KV[K, Vec[V]]) MDict[K, V] {
	d[kv.K] = kv.V
	return d
}
