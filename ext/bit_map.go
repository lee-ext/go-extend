package ext

import "unsafe"

type BytesBitMap struct {
	value []byte
}

func BytesBitMap_(value []byte) BytesBitMap {
	return BytesBitMap{value}
}

func (bm *BytesBitMap) Value() []byte {
	return bm.value
}

// Set 把T的第index位置为b
func (bm *BytesBitMap) Set(index int, b bool) {
	outIndex := index / 8
	if outIndex >= len(bm.value) {
		value := make([]byte, outIndex+1)
		copy(value, bm.value)
		bm.value = value
	}
	v := bm.value[outIndex]
	inIndex := index % 8
	if b {
		v |= byte(1) << inIndex
	} else {
		v &= ^(byte(1) << inIndex)
	}
	bm.value[outIndex] = v
}

// Get 获取T中的第index位置的值
func (bm *BytesBitMap) Get(index int) bool {
	if index >= len(bm.value)*8 {
		return false
	}
	return bm.get(index)
}

func (bm *BytesBitMap) Count() int {
	count := 0
	len_ := len(bm.value) * 8
	for i := range len_ {
		if bm.get(i) {
			count += 1
		}
	}
	return count
}

func (bm *BytesBitMap) Len() int {
	return len(bm.value) * 8
}

func (bm *BytesBitMap) get(index int) bool {
	return (bm.value[index/8]>>(index%8))&1 == 1
}

type BitMap[T Integer] struct {
	value T
}

func BitMap_[T Integer](value T) BitMap[T] {
	return BitMap[T]{value}
}

func (bm *BitMap[T]) Value() T {
	return bm.value
}

// Set 把T的第index位置为b
func (bm *BitMap[T]) Set(index int, b bool) {
	v := bm.value
	if b {
		v |= T(1) << index
	} else {
		v &= ^(T(1) << index)
	}
	bm.value = v
}

// Get 获取T中的第index位置的值
func (bm *BitMap[T]) Get(index int) bool {
	if index >= int(unsafe.Sizeof(bm.value))*8 {
		return false
	}
	return bm.get(index)
}

func (bm *BitMap[T]) Count() int {
	count := 0
	len_ := int(unsafe.Sizeof(bm.value)) * 8
	for i := range len_ {
		if bm.get(i) {
			count += 1
		}
	}
	return count
}

func (bm *BitMap[T]) get(index int) bool {
	return (bm.value>>index)&1 == 1
}

type Bytes2BitMap struct {
	bm BytesBitMap
}

func Bytes2BitMap_(value []byte) Bytes2BitMap {
	return Bytes2BitMap{BytesBitMap{value}}
}

func (bm *Bytes2BitMap) Value() []byte {
	return bm.bm.value
}

func (bm *Bytes2BitMap) Set(index int, v uint8) {
	if v < 1<<2 {
		index *= 2
		for i := range 2 {
			bm.bm.Set(index+i, (v>>i)&1 == 1)
		}
	}
}

func (bm *Bytes2BitMap) Get(index int) uint8 {
	v := 0
	index *= 2
	for i := range 2 {
		if bm.bm.Get(index + i) {
			v += 1 << i
		}
	}
	return uint8(v)
}

func (bm *Bytes2BitMap) Len() int {
	return bm.bm.Len() / 2
}
