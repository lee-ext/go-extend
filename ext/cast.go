package ext

import (
	"unsafe"
)

// UnsafeCast Be cautious when performing binary reinterpretation on value types.
func UnsafeCast[R, T any](t T) R {
	return *(*R)(unsafe.Pointer(&t))
}

// BytesCastStr Avoid unsafe conversions that cause memory copies.
func BytesCastStr(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

// StrCastBytes Avoid unsafe conversions that cause memory copies.
func StrCastBytes(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

// BytesCastNumber Big-endian conversion
func BytesCastNumber[T Number](bytes []byte) T {
	t := *new(T)
	size := unsafe.Sizeof(t)
	switch size {
	case 1:
		t = UnsafeCast[T](bytes[0])
	case 2:
		t = UnsafeCast[T](_BE.Uint16(bytes))
	case 4:
		t = UnsafeCast[T](_BE.Uint32(bytes))
	case 8:
		t = UnsafeCast[T](_BE.Uint64(bytes))
	}
	return t
}

// BytesCastNumberLe Little-endian conversion
func BytesCastNumberLe[T Number](bytes []byte) T {
	t := *new(T)
	size := unsafe.Sizeof(t)
	switch size {
	case 1:
		t = UnsafeCast[T](bytes[0])
	case 2:
		t = UnsafeCast[T](_LE.Uint16(bytes))
	case 4:
		t = UnsafeCast[T](_LE.Uint32(bytes))
	case 8:
		t = UnsafeCast[T](_LE.Uint64(bytes))
	}
	return t
}
