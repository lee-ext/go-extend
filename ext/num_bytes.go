package ext

import (
	"fmt"
	"slices"
	"unsafe"
)

func NumberToBytes[T Number](t T) []byte {
	size := int(unsafe.Sizeof(*new(T)))
	u := UnsafeCast[uint64](t)
	bytes := make([]byte, size)
	for i := range size {
		bytes[i] |= byte(u >> (i * 8))
	}
	return bytes
}

func BytesToNumber[T Number](bytes []byte) T {
	size := int(unsafe.Sizeof(*new(T)))
	if len(bytes) < size {
		return T(0)
	}
	t := uint64(0)
	for i, b := range bytes[0:size] {
		t |= uint64(b) << (i * 8)
	}
	return UnsafeCast[T](t)
}

func NumbersToBytes[T Number](numbers Vec[T]) []byte {
	size := int(unsafe.Sizeof(*new(T)))
	numbers = numbers.Clone()
	bytes := UnsafeCast[[]byte](numbers)
	ptr := (*[3]int)(unsafe.Pointer(&bytes))
	ptr[1], ptr[2] = ptr[1]*size, ptr[2]*size
	return bytes
}

func BytesToNumbers[T Number](bytes []byte) Vec[T] {
	size := int(unsafe.Sizeof(*new(T)))
	if len(bytes)%size > 0 {
		panic(fmt.Errorf("bad bytes len: %v", len(bytes)))
	}
	bytes = slices.Clone(bytes)
	numbers := UnsafeCast[Vec[T]](bytes)
	ptr := (*[3]int)(unsafe.Pointer(&numbers))
	ptr[1], ptr[2] = ptr[1]/size, ptr[2]/size
	return numbers
}
