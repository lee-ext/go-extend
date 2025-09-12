package test

import (
	"fmt"
	"testing"

	"github.com/lee-ext/go-extend/ext"
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
}
