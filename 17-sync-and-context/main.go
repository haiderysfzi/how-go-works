// Lesson 17: Mutexes, atomics, Once, and context for cancellation.
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type SafeCounter struct {
	mu sync.Mutex
	m  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key]++
}

func slowOp(ctx context.Context) error {
	select {
	case <-time.After(500 * time.Millisecond):
		return nil
	case <-ctx.Done(): // cancelled or timed out
		return ctx.Err()
	}
}

type ctxKey string

func main() {
	// Mutex protects shared state. Without it, concurrent map writes crash the program.
	c := SafeCounter{m: map[string]int{}}
	var wg sync.WaitGroup
	for range 1000 {
		wg.Add(1)
		go func() { defer wg.Done(); c.Inc("hits") }()
	}
	wg.Wait()
	fmt.Println(c.m["hits"])

	// RWMutex: many readers OR one writer. Atomics: lock-free single values.
	var rw sync.RWMutex
	rw.RLock()
	rw.RUnlock()
	var n atomic.Int32
	n.Add(5)
	fmt.Println(n.Load(), n.CompareAndSwap(5, 9), n.Load())

	// Once: run initialisation exactly once, even with many goroutines.
	var once sync.Once
	for range 3 {
		once.Do(func() { fmt.Println("init once") })
	}
	lazy := sync.OnceValue(func() int { fmt.Println("computing"); return 42 })
	fmt.Println(lazy(), lazy())

	// sync.Pool reuses temporary objects to reduce GC pressure.
	pool := sync.Pool{New: func() any { return new([]byte) }}
	b := pool.Get().(*[]byte)
	pool.Put(b)

	// context: deadlines, cancellation and request-scoped values, passed as the FIRST param.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()           // always call cancel to release resources
	fmt.Println(slowOp(ctx)) // context deadline exceeded

	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() { time.Sleep(10 * time.Millisecond); cancel2() }()
	fmt.Println(slowOp(ctx2)) // context canceled

	ctx3 := context.WithValue(context.Background(), ctxKey("user"), "ada")
	fmt.Println(ctx3.Value(ctxKey("user")))

	// Condition for choosing: use channels to pass ownership / coordinate,
	// use mutexes to protect a data structure's internal state.
}
