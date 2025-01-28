package test

import (
	"fmt"
	"github.com/lee157953/go-extend/ext"
	"testing"
)

func TestBytes(t *testing.T) {
	bytes := ext.Bytes_(1024)
	bytes.WriteInt32(0, 520)
	bytes.WriteInt64Le(4, 123)
	bytes.WriteFloat64(12, 2.25)
	str := "hello world"
	bytes.WriteString(20, str)
	fmt.Println(bytes.ReadInt32(0))
	fmt.Println(bytes.ReadInt64Le(4))
	fmt.Println(bytes.ReadFloat64(12))
	fmt.Println(bytes.ReadString(20, len(str), false))
}
