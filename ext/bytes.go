package ext

import (
	"encoding/binary"
)

var (
	_BE = binary.BigEndian
	_LE = binary.LittleEndian
)

type Bytes []byte

func Bytes_(len int) Bytes {
	return make([]byte, len)
}

func (b Bytes) Len() int {
	return len(b)
}

func (b Bytes) Empty() bool {
	return len(b) == 0
}

func (b Bytes) Cap() int {
	return cap(b)
}

func (b Bytes) ForEach(fn func(byte)) {
	for _, e := range b {
		fn(e)
	}
}

func (b Bytes) ForEachWhile(fn func(byte) bool) {
	for _, e := range b {
		if !fn(e) {
			return
		}
	}
}

func (b Bytes) ReadBytes(offset, len int) Bytes {
	return b[offset : offset+len]
}

func (b Bytes) ReadString(offset, len int, copy bool) string {
	bytes := b.ReadBytes(offset, len)
	if copy {
		return string(bytes)
	} else {
		return BytesCastStr(bytes)
	}
}

func (b Bytes) WriteBytes(offset int, bytes Bytes) {
	copy(b[offset:offset+bytes.Len()], bytes)
}

func (b Bytes) WriteString(offset int, str string) {
	copy(b[offset:offset+len(str)], str)
}

func (b Bytes) ReadInt8(offset int) int8 {
	return UnsafeCast[int8](b[offset])
}

func (b Bytes) ReadInt16(offset int) int16 {
	return UnsafeCast[int16](_BE.Uint16(b[offset:]))
}

func (b Bytes) ReadInt32(offset int) int32 {
	return UnsafeCast[int32](_BE.Uint32(b[offset:]))
}

func (b Bytes) ReadInt64(offset int) int64 {
	return UnsafeCast[int64](_BE.Uint64(b[offset:]))
}

func (b Bytes) ReadUInt8(offset int) uint8 {
	return b[offset]
}

func (b Bytes) ReadUInt16(offset int) uint16 {
	return _BE.Uint16(b[offset:])
}

func (b Bytes) ReadUInt32(offset int) uint32 {
	return _BE.Uint32(b[offset:])
}

func (b Bytes) ReadUInt64(offset int) uint64 {
	return _BE.Uint64(b[offset:])
}

func (b Bytes) ReadFloat32(offset int) float32 {
	return UnsafeCast[float32](_BE.Uint32(b[offset:]))
}

func (b Bytes) ReadFloat64(offset int) float64 {
	return UnsafeCast[float64](_BE.Uint64(b[offset:]))
}

func (b Bytes) WriteInt8(offset int, value int8) {
	b[offset] = UnsafeCast[uint8](value)
}

func (b Bytes) WriteInt16(offset int, value int16) {
	_BE.PutUint16(b[offset:], UnsafeCast[uint16](value))
}

func (b Bytes) WriteInt32(offset int, value int32) {
	_BE.PutUint32(b[offset:], UnsafeCast[uint32](value))
}

func (b Bytes) WriteInt64(offset int, value int64) {
	_BE.PutUint64(b[offset:], UnsafeCast[uint64](value))
}

func (b Bytes) WriteUInt8(offset int, value uint8) {
	b[offset] = value
}

func (b Bytes) WriteUInt16(offset int, value uint16) {
	_BE.PutUint16(b[offset:], value)
}

func (b Bytes) WriteUInt32(offset int, value uint32) {
	_BE.PutUint32(b[offset:], value)
}

func (b Bytes) WriteUInt64(offset int, value uint64) {
	_BE.PutUint64(b[offset:], value)
}

func (b Bytes) WriteFloat32(offset int, value float32) {
	_BE.PutUint32(b[offset:], UnsafeCast[uint32](value))
}

func (b Bytes) WriteFloat64(offset int, value float64) {
	_BE.PutUint64(b[offset:], UnsafeCast[uint64](value))
}

func (b Bytes) ReadInt16Le(offset int) int16 {
	return UnsafeCast[int16](_LE.Uint16(b[offset:]))
}

func (b Bytes) ReadInt32Le(offset int) int32 {
	return UnsafeCast[int32](_LE.Uint32(b[offset:]))
}

func (b Bytes) ReadInt64Le(offset int) int64 {
	return UnsafeCast[int64](_LE.Uint64(b[offset:]))
}

func (b Bytes) ReadUInt16Le(offset int) uint16 {
	return _LE.Uint16(b[offset:])
}

func (b Bytes) ReadUInt32Le(offset int) uint32 {
	return _LE.Uint32(b[offset:])
}

func (b Bytes) ReadUInt64Le(offset int) uint64 {
	return _LE.Uint64(b[offset:])
}

func (b Bytes) ReadFloat32Le(offset int) float32 {
	return UnsafeCast[float32](_LE.Uint32(b[offset:]))
}

func (b Bytes) ReadFloat64Le(offset int) float64 {
	return UnsafeCast[float64](_LE.Uint64(b[offset:]))
}

func (b Bytes) WriteInt16Le(offset int, value int16) {
	_LE.PutUint16(b[offset:], UnsafeCast[uint16](value))
}

func (b Bytes) WriteInt32Le(offset int, value int32) {
	_LE.PutUint32(b[offset:], UnsafeCast[uint32](value))
}

func (b Bytes) WriteInt64Le(offset int, value int64) {
	_LE.PutUint64(b[offset:], UnsafeCast[uint64](value))
}

func (b Bytes) WriteUInt16Le(offset int, value uint16) {
	_LE.PutUint16(b[offset:], value)
}

func (b Bytes) WriteUInt32Le(offset int, value uint32) {
	_LE.PutUint32(b[offset:], value)
}

func (b Bytes) WriteUInt64Le(offset int, value uint64) {
	_LE.PutUint64(b[offset:], value)
}

func (b Bytes) WriteFloat32Le(offset int, value float32) {
	_LE.PutUint32(b[offset:], UnsafeCast[uint32](value))
}

func (b Bytes) WriteFloat64Le(offset int, value float64) {
	_LE.PutUint64(b[offset:], UnsafeCast[uint64](value))
}
