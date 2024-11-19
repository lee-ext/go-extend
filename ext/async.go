package ext

import (
	"fmt"
	"sync"
	"time"
)

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

type _FuturePin[T any] struct {
	waiter sync.WaitGroup
	result CuRes[T]
}

type Future[T any] struct {
	*_FuturePin[T]
}

func (f Future[T]) Await() CuRes[T] {
	f.waiter.Wait()
	return f.result
}

func Async[T any](fn func() T) Future[T] {
	f := Future[T]{&_FuturePin[T]{}}
	f.waiter.Add(1)
	go func() {
		defer func() {
			switch e := recover().(type) {
			case error:
				f.result = CuErr[T](e)
			case nil:
				break
			default:
				f.result = CuErr[T](fmt.Errorf("unknown error: %#v", e))
			}
			f.waiter.Done()
		}()
		f.result = CuOk(fn())
	}()
	return f
}
