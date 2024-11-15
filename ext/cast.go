package ext

import (
	"encoding/binary"
	"unsafe"
)

// UnsafeCast 对值类型做二进制重新解释 谨慎使用
func UnsafeCast[R, T any](t T) R {
	return *(*R)(unsafe.Pointer(&t))
}

// BytesCastStr 为了减少内存copy 使用的不安全转换 谨慎使用 传入的bytes不允许再做修改
func BytesCastStr[T ~[]byte](bytes T) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

// StrCastBytes 为了减少内存copy 使用的不安全转换 谨慎使用
func StrCastBytes[T ~string](str T) []byte {
	return unsafe.Slice(unsafe.StringData(string(str)), len(str))
}

// BytesCastNumber 二进制转数字
func BytesCastNumber[T Number](bytes []byte) T {
	t := *new(T)
	size := unsafe.Sizeof(t)
	switch size {
	case 1:
		t = UnsafeCast[T](bytes[0])
	case 2:
		t = UnsafeCast[T](binary.BigEndian.Uint16(bytes))
	case 4:
		t = UnsafeCast[T](binary.BigEndian.Uint32(bytes))
	case 8:
		t = UnsafeCast[T](binary.BigEndian.Uint64(bytes))
	}
	return t
}

// BytesCastNumberLe 二进制转数字 小端
func BytesCastNumberLe[T Number](bytes []byte) T {
	t := *new(T)
	size := unsafe.Sizeof(t)
	switch size {
	case 1:
		t = UnsafeCast[T](bytes[0])
	case 2:
		t = UnsafeCast[T](binary.LittleEndian.Uint16(bytes))
	case 4:
		t = UnsafeCast[T](binary.LittleEndian.Uint32(bytes))
	case 8:
		t = UnsafeCast[T](binary.LittleEndian.Uint64(bytes))
	}
	return t
}
