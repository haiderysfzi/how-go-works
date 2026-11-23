// Lesson 17: Mutexes, atomics, Once, and context for cancellation.
//
// WHAT YOU WILL LEARN
//   - What a data race is and how a Mutex prevents it
//   - RWMutex, atomics, sync.Once and sync.Pool
//   - `context`: how to cancel work, set deadlines and carry request-scoped values
//   - When to use channels vs mutexes
//
// BACKGROUND: when two goroutines access the SAME variable at the same time and at
// least one writes, you have a DATA RACE, and the result is unpredictable. Example:
// `counter++` compiles to read -> add -> write. If two goroutines both read 5, both
// compute 6, and both write 6, one increment is lost. The tools below prevent this.
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// A MUTEX (mutual exclusion lock) lets only ONE goroutine at a time into a critical
// section. Lock() waits until the lock is free and takes it; Unlock() releases it.
// Convention: put the mutex right above the data it protects in the same struct.
type SafeCounter struct {
	mu sync.Mutex
	m  map[string]int
}

// `defer c.mu.Unlock()` guarantees the unlock even if the code in between panics or
// returns early. Forgetting to unlock leaves every other goroutine stuck forever.
// The method needs a POINTER receiver: copying a Mutex copies its state and breaks it.
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key]++
}

// CONTEXT USAGE: slowOp pretends to take 500ms. `select` (lesson 16) waits for
// whichever happens first: the work finishing, or ctx.Done(), a channel that is
// CLOSED when the context is cancelled or its deadline passes. ctx.Err() then says
// why: context.Canceled or context.DeadlineExceeded.
func slowOp(ctx context.Context) error {
	select {
	case <-time.After(500 * time.Millisecond):
		return nil
	case <-ctx.Done(): // cancelled or timed out
		return ctx.Err()
	}
}

// Context keys should have their own unexported type so that keys from different
// packages can never collide (a plain "user" string could clash with another package's).
type ctxKey string

func main() {
	// 1000 goroutines each add 1 to the same map entry. The Mutex serialises them,
	// so the final answer is exactly 1000. WITHOUT the Mutex the Go runtime detects
	// concurrent map writes and crashes the program on purpose.
	c := SafeCounter{m: map[string]int{}}
	var wg sync.WaitGroup
	for range 1000 {
		wg.Add(1)
		go func() { defer wg.Done(); c.Inc("hits") }()
	}
	wg.Wait()
	fmt.Println(c.m["hits"]) // 1000

	// RWMutex: a lock that allows MANY simultaneous readers (RLock) OR ONE writer (Lock).
	// Worth it when reads vastly outnumber writes. (Shown only to display the API.)
	var rw sync.RWMutex
	rw.RLock()
	rw.RUnlock()

	// ATOMICS: single-variable operations performed as ONE indivisible CPU instruction,
	// with no lock. Faster than a mutex for simple counters and flags.
	// CompareAndSwap(old, new) sets the value to `new` ONLY if it currently equals
	// `old`, and reports whether it did: the building block for lock-free algorithms.
	// (Printf-style arguments are evaluated left to right: Load = 5, swap -> true, Load = 9.)
	var n atomic.Int32
	n.Add(5)
	fmt.Println(n.Load(), n.CompareAndSwap(5, 9), n.Load()) // 5 true 9

	// ONCE: Do runs the function exactly once, even when called many times and from
	// many goroutines at once; later callers wait until the first call has finished.
	// Typical use: lazy initialisation of a shared resource (a config, a connection).
	var once sync.Once
	for range 3 {
		once.Do(func() { fmt.Println("init once") }) // prints one time only
	}
	// OnceValue (Go 1.21+) wraps a function that computes a value: it is evaluated on
	// the first call and the cached result is returned on every later call.
	lazy := sync.OnceValue(func() int { fmt.Println("computing"); return 42 })
	fmt.Println(lazy(), lazy()) // "computing" is printed once, then: 42 42

	// POOL: a cache of reusable temporary objects. Get() returns a recycled one (or
	// makes a new one with New); Put() gives it back. Reusing buffers means fewer
	// allocations and less garbage-collector work (lesson 18). The runtime may drop
	// pooled objects at any time, so never store anything in a Pool you cannot recreate.
	pool := sync.Pool{New: func() any { return new([]byte) }}
	b := pool.Get().(*[]byte) // Get returns `any`, so we assert the real type
	pool.Put(b)

	// CONTEXT carries three things across function and goroutine boundaries:
	// cancellation signals, deadlines, and request-scoped values. Convention: pass it
	// as the FIRST parameter, named ctx. Contexts form a tree: cancelling a parent
	// cancels all its children.
	//
	// WithTimeout creates a child context that cancels itself after 50 ms. `cancel`
	// is a function that releases the timer; always call it (usually with defer) even
	// if the timeout already fired, or resources linger until the parent ends.
	// slowOp wants 500 ms but only has 50, so it returns "context deadline exceeded".
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()           // always call cancel to release resources
	fmt.Println(slowOp(ctx)) // context deadline exceeded

	// WithCancel: no timer; cancellation happens only when WE call cancel2().
	// Another goroutine does that after 10 ms, so slowOp stops early with
	// "context canceled".
	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() { time.Sleep(10 * time.Millisecond); cancel2() }()
	fmt.Println(slowOp(ctx2)) // context canceled

	// WithValue attaches a key/value pair; Value(key) retrieves it, possibly from a
	// parent. Use it for request-scoped data such as a request ID or authenticated
	// user, NOT as a way to pass ordinary function arguments.
	ctx3 := context.WithValue(context.Background(), ctxKey("user"), "ada")
	fmt.Println(ctx3.Value(ctxKey("user"))) // ada

	// CHOOSING BETWEEN TOOLS
	// - Channels: passing OWNERSHIP of data, pipelines, signalling events.
	// - Mutexes:  protecting a data structure's internal state (like the map above).
	// - Atomics:  a single counter or flag.
	// - Context:  cancellation, deadlines and request-scoped values.
}
