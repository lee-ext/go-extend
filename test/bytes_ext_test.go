package test

import (
	"math"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestBytesNumbers(t *testing.T) {
	b := Bytes_(64)
	b.WriteInt32(0, -520)
	if b.ReadInt32(0) != -520 || b.ReadUInt32(0) != uint32(0xfffffdf8) {
		t.Errorf("Int32 = %d/%d", b.ReadInt32(0), b.ReadUInt32(0))
	}
	b.WriteInt64(8, math.MinInt64)
	if b.ReadInt64(8) != math.MinInt64 {
		t.Errorf("Int64 = %d", b.ReadInt64(8))
	}
	b.WriteInt8(16, -7)
	if b.ReadInt8(16) != -7 || b.ReadUInt8(16) != 249 {
		t.Errorf("Int8 = %d/%d", b.ReadInt8(16), b.ReadUInt8(16))
	}
	b.WriteInt16(17, -300)
	if b.ReadInt16(17) != -300 {
		t.Errorf("Int16 = %d", b.ReadInt16(17))
	}
	b.WriteFloat32(19, 1.5)
	if b.ReadFloat32(19) != 1.5 {
		t.Errorf("Float32 = %v", b.ReadFloat32(19))
	}
	b.WriteFloat64(23, 2.25)
	if b.ReadFloat64(23) != 2.25 {
		t.Errorf("Float64 = %v", b.ReadFloat64(23))
	}
	b.WriteUInt16(31, 65535)
	b.WriteUInt64(33, math.MaxUint64)
	if b.ReadUInt16(31) != 65535 || b.ReadUInt64(33) != math.MaxUint64 {
		t.Errorf("UInt = %d/%d", b.ReadUInt16(31), b.ReadUInt64(33))
	}
}

func TestBytesNumbersLE(t *testing.T) {
	b := Bytes_(64)
	b.WriteInt16Le(0, -300)
	b.WriteInt32Le(2, -520)
	b.WriteInt64Le(6, math.MinInt64)
	b.WriteUInt16Le(14, 65535)
	b.WriteUInt32Le(16, math.MaxUint32)
	b.WriteUInt64Le(20, math.MaxUint64)
	b.WriteFloat32Le(28, 1.5)
	b.WriteFloat64Le(32, 2.25)
	if b.ReadInt16Le(0) != -300 || b.ReadInt32Le(2) != -520 || b.ReadInt64Le(6) != math.MinInt64 {
		t.Errorf("LE ints = %d/%d/%d", b.ReadInt16Le(0), b.ReadInt32Le(2), b.ReadInt64Le(6))
	}
	if b.ReadUInt16Le(14) != 65535 || b.ReadUInt32Le(16) != math.MaxUint32 || b.ReadUInt64Le(20) != math.MaxUint64 {
		t.Error("LE uints mismatch")
	}
	if b.ReadFloat32Le(28) != 1.5 || b.ReadFloat64Le(32) != 2.25 {
		t.Error("LE floats mismatch")
	}
	// writing big-endian then reading little-endian must differ
	b.WriteInt32(40, 0x01020304)
	if b.ReadInt32(40) == b.ReadInt32Le(40) {
		t.Error("big and little endian reads should not agree")
	}
	if b.ReadInt32Le(40) != 0x04030201 {
		t.Errorf("ReadInt32Le = %#x", b.ReadInt32Le(40))
	}
}

func TestBytesStrAndSlices(t *testing.T) {
	b := Bytes_(32)
	b.WriteString(0, "hello")
	if s := b.ReadString(0, 5, true); s != "hello" {
		t.Errorf("ReadString copy = %q", s)
	}
	if s := b.ReadString(0, 5, false); s != "hello" {
		t.Errorf("ReadString view = %q", s)
	}
	b.WriteBytes(6, Bytes{'w', 'o', 'r', 'l', 'd'})
	if got := b.ReadBytes(6, 5); string(got) != "world" {
		t.Errorf("ReadBytes = %q", got)
	}
	if b.Len() != 32 || b.Empty() || b.Cap() != 32 {
		t.Errorf("meta len=%d cap=%d", b.Len(), b.Cap())
	}
	if Bytes_(0).Empty() != true || Bytes_(0).Len() != 0 {
		t.Error("empty Bytes meta wrong")
	}
	count := 0
	b.ForEach(func(byte) { count++ })
	if count != 32 {
		t.Errorf("ForEach count = %d", count)
	}
	n := 0
	b.ForEachWhile(func(byte) bool { n++; return n < 3 })
	if n != 3 {
		t.Errorf("ForEachWhile n = %d", n)
	}
}

func TestBitMap(t *testing.T) {
	bm := BitMap_[uint8](0)
	bm.Set(0, true)
	bm.Set(3, true)
	if bm.Value() != 0b1001 || !bm.Get(0) || !bm.Get(3) || bm.Get(1) {
		t.Fatalf("BitMap = %08b", bm.Value())
	}
	if bm.Get(8) || bm.Get(100) {
		t.Error("out of range Get should be false")
	}
	if bm.Count() != 2 {
		t.Errorf("Count = %d", bm.Count())
	}
	bm.Set(0, false)
	if bm.Get(0) || bm.Value() != 0b1000 {
		t.Errorf("clear bit = %08b", bm.Value())
	}
	u64 := BitMap_[uint64](0)
	u64.Set(63, true)
	if !u64.Get(63) || u64.Count() != 1 {
		t.Errorf("uint64 bit63 = %d", u64.Count())
	}
}

func TestBytesBitMap(t *testing.T) {
	bm := BytesBitMap_(make([]byte, 1))
	bm.Set(0, true)
	bm.Set(9, true)
	if !bm.Get(0) || !bm.Get(9) || bm.Get(1) {
		t.Fatalf("BytesBitMap get = %v", bm.Value())
	}
	// Set(9) grows the backing slice to 2 bytes, so Len() is 16 bits.
	if bm.Count() != 2 || bm.Len() != 16 {
		t.Errorf("Count=%d Len=%d", bm.Count(), bm.Len())
	}
	if bm.Get(100) {
		t.Error("out of range Get should be false")
	}
	// Set beyond capacity grows the backing slice
	big := BytesBitMap_(nil)
	big.Set(40, true)
	if !big.Get(40) || big.Count() != 1 || len(big.Value()) < 6 {
		t.Errorf("grow on Set = %v count=%d", big.Value(), big.Count())
	}
	big.Set(40, false)
	if big.Get(40) || big.Count() != 0 {
		t.Errorf("clear grown bit = %d", big.Count())
	}
}

func TestBytes2BitMap(t *testing.T) {
	bm := Bytes2BitMap_(make([]byte, 1))
	bm.Set(0, 3)
	bm.Set(1, 1)
	bm.Set(2, 2)
	if bm.Get(0) != 3 || bm.Get(1) != 1 || bm.Get(2) != 2 {
		t.Fatalf("Bytes2BitMap = %d,%d,%d", bm.Get(0), bm.Get(1), bm.Get(2))
	}
	if bm.Len() != 4 {
		t.Errorf("Len = %d", bm.Len())
	}
	// values >= 1<<2 are ignored by Set
	bm.Set(0, 9)
	if bm.Get(0) != 3 {
		t.Errorf("Set out of range should be ignored, got %d", bm.Get(0))
	}
}

func TestNumberBytesCast(t *testing.T) {
	if got := NumberToBytes(int32(0x01020304)); len(got) != 4 {
		t.Fatalf("NumberToBytes len = %d", len(got))
	}
	for _, v := range []int64{0, 1, -1, math.MaxInt64, math.MinInt64, 12345} {
		if got := BytesToNumber[int64](NumberToBytes(v)); got != v {
			t.Errorf("int64 round trip %d -> %d", v, got)
		}
	}
	if got := BytesToNumber[int32]([]byte{1}); got != 0 {
		t.Errorf("short bytes should yield 0, got %d", got)
	}
	f := 3.5
	if got := BytesToNumber[float64](NumberToBytes(f)); got != f {
		t.Errorf("float64 round trip %v -> %v", f, got)
	}

	nums := VecOf[uint16](1, 2, 3)
	bs := NumbersToBytes(nums)
	if len(bs) != 6 {
		t.Fatalf("NumbersToBytes len = %d", len(bs))
	}
	back := BytesToNumbers[uint16](bs)
	if len(back) != 3 || back[0] != 1 || back[2] != 3 {
		t.Errorf("BytesToNumbers = %v", back)
	}
	if !panics(func() { BytesToNumbers[uint16]([]byte{1, 2, 3}) }) {
		t.Error("BytesToNumbers should panic on misaligned input")
	}
}

func TestUnsafeCastAndStr(t *testing.T) {
	if UnsafeCast[uint8](int8(-1)) != 255 {
		t.Errorf("UnsafeCast int8->uint8 = %d", UnsafeCast[uint8](int8(-1)))
	}
	if UnsafeCast[int32](uint32(0xffffffff)) != -1 {
		t.Errorf("UnsafeCast uint32->int32 = %d", UnsafeCast[int32](uint32(0xffffffff)))
	}
	if BytesCastNumber[int32]([]byte{0x01, 0x02, 0x03, 0x04}) != 0x01020304 {
		t.Errorf("BytesCastNumber BE = %d", BytesCastNumber[int32]([]byte{1, 2, 3, 4}))
	}
	if BytesCastNumberLe[int32]([]byte{0x01, 0x02, 0x03, 0x04}) != 0x04030201 {
		t.Errorf("BytesCastNumberLe = %d", BytesCastNumberLe[int32]([]byte{1, 2, 3, 4}))
	}
	if BytesCastNumber[int8]([]byte{0xff}) != -1 {
		t.Errorf("BytesCastNumber int8 = %d", BytesCastNumber[int8]([]byte{0xff}))
	}
	if BytesCastNumberLe[uint16]([]byte{0x01, 0x02}) != 0x0201 {
		t.Errorf("BytesCastNumberLe uint16 = %d", BytesCastNumberLe[uint16]([]byte{1, 2}))
	}
	s := "round trip"
	if BytesCastStr(StrCastBytes(s)) != s {
		t.Error("string/bytes cast round trip failed")
	}
}
