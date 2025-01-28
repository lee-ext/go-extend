package ext

import (
	"fmt"
	"time"
)

// DeferFn_ example
func DeferFn_(r any) {
	switch e := r.(type) {
	case nil:
		break
	case error:
		fmt.Printf("Error: %v\n", e.Error())
	default:
		fmt.Printf("Unknown error: %#v\n", e)
	}
}

// Launch a new function, deferFn used for recover processing,
// panic does not cause a crash
func Launch(fn func(), deferFn func(any)) {
	go func() {
		defer func() {
			deferFn(recover())
		}()
		fn()
	}()
}

// LaunchSturdy same of Launch, but this will wake itself up
func LaunchSturdy(fn func(), deferFn func(any)) {
	go func() {
		defer func() {
			deferFn(recover())
			time.Sleep(time.Second)
			LaunchSturdy(fn, deferFn)
		}()
		fn()
	}()
}
