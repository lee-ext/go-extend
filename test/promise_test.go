package test

import (
	"sync"
	"testing"
	"time"

	. "github.com/lee-ext/go-extend/ext"
)

func TestPromiseBasic(t *testing.T) {
	p := Promise_[int]()
	if !p.Pending() || p.Completed() {
		t.Fatal("new promise should be pending")
	}
	if p.TryGet().IsSome() {
		t.Error("pending TryGet should be None")
	}
	if !p.Complete(7) {
		t.Error("first Complete should succeed")
	}
	if p.Pending() || !p.Completed() {
		t.Error("promise should be completed")
	}
	if p.Await() != 7 {
		t.Errorf("Await = %d", p.Await())
	}
	if got := p.TryGet(); !got.IsSome() || got.Get() != 7 {
		t.Errorf("TryGet = %v", got)
	}
	if p.Complete(99) {
		t.Error("second Complete should fail")
	}
	if p.Await() != 7 {
		t.Errorf("value must not change, got %d", p.Await())
	}
}

func TestPromiseConcurrent(t *testing.T) {
	p := Promise_[string]()
	var wg sync.WaitGroup
	var winners int64
	var mu sync.Mutex
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if p.Complete("done") {
				mu.Lock()
				winners++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if winners != 1 {
		t.Errorf("exactly one Complete should win, got %d", winners)
	}
	if p.Await() != "done" {
		t.Errorf("Await = %q", p.Await())
	}
}

func TestPromiseAwaitBlocks(t *testing.T) {
	p := Promise_[int]()
	go func() {
		time.Sleep(20 * time.Millisecond)
		p.Complete(5)
	}()
	done := make(chan int, 1)
	go func() { done <- p.Await() }()
	select {
	case v := <-done:
		if v != 5 {
			t.Errorf("Await = %d", v)
		}
	case <-time.After(time.Second):
		t.Fatal("Await did not return after Complete")
	}
}

func TestPromiseCompleteAfter(t *testing.T) {
	p := Promise_[int]()
	timer := p.CompleteAfter(10*time.Millisecond, func() int { return 42 })
	defer timer.Stop()
	if v := p.Await(); v != 42 {
		t.Errorf("CompleteAfter = %d", v)
	}
	if !p.Completed() {
		t.Error("should be completed")
	}
	// nil fn must panic immediately
	if !panics(func() { Promise_[int]().CompleteAfter(time.Second, nil) }) {
		t.Error("CompleteAfter(nil) should panic")
	}
	// late CompleteAfter on an already completed promise is a no-op
	q := Promise_[int]()
	q.Complete(1)
	q.CompleteAfter(time.Millisecond, func() int { return 2 }).Stop()
	time.Sleep(5 * time.Millisecond)
	if q.Await() != 1 {
		t.Errorf("late CompleteAfter changed value: %d", q.Await())
	}
}

func TestActorLaunch(t *testing.T) {
	a := Actor_(16, DeferFn_)
	defer a.Close()
	results := make([]int, 5)
	done := Promise_[Unit]()
	for i := range 5 {
		i := i
		a.Launch(func() {
			results[i] = i * 2
			if i == 4 {
				done.Complete(Unit{})
			}
		})
	}
	done.Await()
	for i, v := range results {
		if v != i*2 {
			t.Errorf("results[%d] = %d", i, v)
		}
	}
}

func TestActorAwait(t *testing.T) {
	a := Actor_(8, DeferFn_)
	defer a.Close()
	p := a.Await(func() string { return "hello actor" })
	if got := p.Await(); got != "hello actor" {
		t.Errorf("Actor.Await = %q", got)
	}
}

func TestActorSerialOrder(t *testing.T) {
	a := Actor_(64, DeferFn_)
	defer a.Close()
	var order []int
	var mu sync.Mutex
	done := Promise_[Unit]()
	for i := range 20 {
		i := i
		a.Launch(func() {
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
			if i == 19 {
				done.Complete(Unit{})
			}
		})
	}
	done.Await()
	for i, v := range order {
		if v != i {
			t.Fatalf("actor must run tasks serially in order, got %v", order)
		}
	}
}

func TestActorRecoversPanic(t *testing.T) {
	var mu sync.Mutex
	panics := 0
	a := Actor_(8, func(any) {
		mu.Lock()
		panics++
		mu.Unlock()
	})
	defer a.Close()
	a.Launch(func() { panic("boom") })
	// give the recover handler time to run
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		n := panics
		mu.Unlock()
		if n >= 1 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	mu.Lock()
	defer mu.Unlock()
	if panics < 1 {
		t.Errorf("deferFn should observe the panic, got %d", panics)
	}
}

func TestLaunch(t *testing.T) {
	done := Promise_[int]()
	Launch(func() { done.Complete(3) }, DeferFn_)
	if v := done.Await(); v != 3 {
		t.Errorf("Launch result = %d", v)
	}
	recovered := Promise_[any]()
	Launch(func() { panic("boom") }, func(r any) {
		if !recovered.Completed() {
			recovered.Complete(r)
		}
	})
	if r := recovered.Await(); r == nil {
		t.Error("Launch deferFn should receive the panic value")
	}
}

func TestPool(t *testing.T) {
	created := 0
	var mu sync.Mutex
	p := Pool_(func() *bytes_buf {
		mu.Lock()
		created++
		mu.Unlock()
		return &bytes_buf{}
	})
	b := p.Get()
	if b == nil {
		t.Fatal("Pool.Get returned nil")
	}
	b.n = 5
	p.Put(b)
	if got := p.Get(); got == nil {
		t.Fatal("Pool.Get after Put returned nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if created < 1 {
		t.Errorf("New should run at least once, created=%d", created)
	}
}

type bytes_buf struct {
	n int
}
