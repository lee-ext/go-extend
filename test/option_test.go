package test

import (
	"fmt"
	"github.com/lee157953/go-extend/ext"
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	opt := ext.Opt_(0, false)
	func() {
		defer func() {
			switch r := recover().(type) {
			case nil:
				break
			case error:
				fmt.Println(r.Error())
			default:
				fmt.Println(r)
			}
		}()
		fmt.Println(opt.Get())
	}()
	fmt.Println(opt.Get_())
	fmt.Println(opt.GetOr(1))
	fmt.Println(opt.GetElse(func() int { return 2 }))

	sum := int64(0)
	t0 := time.Now().UnixMilli()
	for i := int64(0); i <= 1000000000; i++ {
		sum += i
	}
	t1 := time.Now().UnixMilli()
	println("sum:", sum)
	println("duration:", t1-t0)
}
