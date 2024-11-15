package ext

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

type OptU[T uint | uint8 | uint16 | uint32 | uint64] struct {
	u T
}

func OptU_[T uint | uint8 | uint16 | uint32 | uint64](u T, b bool) OptU[T] {
	o := OptU[T]{}
	if b {
		o.Set(u)
	}
	return o
}

func (o OptU[T]) IsSome() bool {
	return o.u > 0
}

func (o OptU[T]) IsNone() bool {
	return !o.IsSome()
}

func (o *OptU[T]) Set(u T) {
	(*o).u = u + 1
}

func (o OptU[T]) get() T {
	return o.u - 1
}

func (o OptU[T]) Get() T {
	if o.IsSome() {
		return o.get()
	}
	panic(errors.New("option is none"))
}

// Get_ 获取值 如果为none 则会返回初始值
func (o OptU[T]) Get_() T {
	if o.IsSome() {
		return o.get()
	}
	return T(0)
}

func (o OptU[T]) GetOr(t T) T {
	if o.IsSome() {
		return o.get()
	}
	return t
}

func (o OptU[T]) GetElse(f func() T) T {
	if o.IsSome() {
		return o.get()
	}
	return f()
}

// MarshalJSON returns m as the JSON encoding of m.
func (o OptU[T]) MarshalJSON() ([]byte, error) {
	if o.IsSome() {
		return []byte(intToStr[T](o.get())), nil
	}
	return []byte("null"), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (o *OptU[T]) UnmarshalJSON(data []byte) error {
	if o == nil {
		return errors.New("UnmarshalJSON on nil pointer")
	}
	str := string(data)
	if str != "null" {
		u, err := strToInt[T](str)
		if err != nil {
			return err
		}
		(*o).u = u + 1
	}
	return nil
}

func (o OptU[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("some(%v)", o.get())
	}
	return "none"
}

type OptI[T int | int8 | int16 | int32 | int64] struct {
	i T
}

func OptI_[T int | int8 | int16 | int32 | int64](i T, b bool) OptI[T] {
	o := OptI[T]{}
	if b {
		o.Set(i)
	}
	return o
}

func (o OptI[T]) IsSome() bool {
	return o.i != 0
}

func (o OptI[T]) IsNone() bool {
	return !o.IsSome()
}

func (o *OptI[T]) Set(i T) {
	if i <= 0 {
		i -= 1
	}
	(*o).i = i
}

func (o OptI[T]) get() T {
	if o.i < 0 {
		o.i += 1
	}
	return o.i
}

func (o OptI[T]) Get() T {
	if o.IsSome() {
		return o.get()
	}
	panic(errors.New("option is none"))
}

// Get_ 获取值 如果为none 则会返回初始值
func (o OptI[T]) Get_() T {
	if o.IsSome() {
		return o.get()
	}
	return T(0)
}

func (o OptI[T]) GetOr(t T) T {
	if o.IsSome() {
		return o.get()
	}
	return t
}

func (o OptI[T]) GetElse(f func() T) T {
	if o.IsSome() {
		return o.get()
	}
	return f()
}

// MarshalJSON returns m as the JSON encoding of m.
func (o OptI[T]) MarshalJSON() ([]byte, error) {
	if o.IsSome() {
		return []byte(intToStr[T](o.get())), nil
	}
	return []byte("null"), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (o *OptI[T]) UnmarshalJSON(data []byte) error {
	if o == nil {
		return errors.New("UnmarshalJSON on nil pointer")
	}
	str := string(data)
	if str != "null" {
		i, err := strToInt[T](str)
		if err != nil {
			return err
		}
		o.Set(i)
	}
	return nil
}

func (o OptI[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("some(%v)", o.get())
	}
	return "none"
}

type OptF[T float32 | float64] struct {
	f T
}

func OptF_[T float32 | float64](f T, b bool) OptF[T] {
	o := OptF[T]{}
	if b {
		o.Set(f)
	}
	return o
}

func (o OptF[T]) Opt() Opt[T] {
	return Opt_(o.Get(), o.IsSome())
}

func (o OptF[T]) IsSome() bool {
	if o.f != 0 {
		return true
	} else {
		switch unsafe.Sizeof(o.f) {
		case 4:
			return math.Float32bits(float32(o.f))&(1<<31) != 0
		case 8:
			return math.Float64bits(float64(o.f))&(1<<63) != 0
		default:
			return false
		}
	}
}

func (o OptF[T]) IsNone() bool {
	return !o.IsSome()
}

func (o *OptF[T]) Set(f T) {
	if f == 0 {
		f = -f
	}
	(*o).f = f
}

func (o OptF[T]) get() T {
	if o.f == 0 {
		o.f = -o.f
	}
	return o.f
}

func (o OptF[T]) Get() T {
	if o.IsSome() {
		return o.get()
	}
	panic(errors.New("option is none"))
}

// Get_ 获取值 如果为none 则会返回初始值
func (o OptF[T]) Get_() T {
	if o.IsSome() {
		return o.get()
	}
	return T(0)
}

func (o OptF[T]) GetOr(t T) T {
	if o.IsSome() {
		return o.get()
	}
	return t
}

func (o OptF[T]) GetElse(f func() T) T {
	if o.IsSome() {
		return o.f
	}
	return f()
}

// MarshalJSON returns m as the JSON encoding of m.
func (o OptF[T]) MarshalJSON() ([]byte, error) {
	if o.IsSome() {
		return []byte(floatToStr(o.f)), nil
	}
	return []byte("null"), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (o *OptF[T]) UnmarshalJSON(data []byte) error {
	if o == nil {
		return errors.New("UnmarshalJSON on nil pointer")
	}
	str := string(data)
	if str != "null" {
		f, err := strToFloat[T](str)
		if err != nil {
			return err
		}
		o.Set(f)
	}
	return nil
}

func (o OptF[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("some(%v)", o.get())
	}
	return "none"
}

func intToStr[T Integer](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func strToInt[T Integer](s string) (T, error) {
	size := int(unsafe.Sizeof(*new(T)))
	i, err := strconv.ParseInt(s, 10, size*8)
	return T(i), err
}

func floatToStr[T Float](f T) string {
	size := int(unsafe.Sizeof(*new(T)))
	return strconv.FormatFloat(float64(f), 'f', -1, size*8)
}

func strToFloat[T Float](s string) (T, error) {
	size := int(unsafe.Sizeof(*new(T)))
	f, err := strconv.ParseFloat(s, size*8)
	return T(f), err
}
