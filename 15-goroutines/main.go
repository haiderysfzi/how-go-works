// Lesson 15: Goroutines: cheap, runtime-managed concurrent functions.
//
// WHAT YOU WILL LEARN
//   - What a goroutine is and how it differs from an OS thread
//   - How to start one, and how to wait for several to finish (sync.WaitGroup)
//   - Why Go can run 100,000 of them
//   - Data races, atomics, and how the per-iteration loop variable helps
//
// BACKGROUND: normally a program runs one line after another (one "thread of
// execution"). A THREAD is something the operating system schedules onto a CPU core;
// OS threads are heavyweight (often 1-8 MB of stack each, and switching between
// them goes through the kernel). A GOROUTINE is a function running independently of
// its caller, managed by the Go RUNTIME instead of the OS. It starts with about 2 KB
// of stack that grows when needed, and the runtime multiplexes thousands of goroutines
// onto a handful of OS threads (lesson 19 shows how).
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// A sync.WaitGroup is a counter. wg.Add(1) increments it, wg.Done() decrements it
// (it is shorthand for Add(-1)), and wg.Wait() blocks until the counter returns to zero.
// It is passed as a POINTER (*sync.WaitGroup): a copy would have its own separate
// counter and Wait would never notice the Done calls.
// `defer wg.Done()` guarantees the counter is decremented even if the function panics.
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(id) * time.Millisecond) // pretend to work; staggered so output order is stable
	fmt.Println("worker", id, "done")
}

func main() {
	// `go f(x)` starts f(x) in a NEW goroutine and returns IMMEDIATELY; the caller does
	// not wait. The runtime decides when and on which thread the new goroutine runs.
	// Hence the WaitGroup: if main ended first, the process would exit and kill the workers.
	// Add(1) is called BEFORE `go`, in the parent, so Wait cannot observe a zero counter
	// before the goroutine has had a chance to start.
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait() // block until the counter hits zero

	// GOROUTINES ARE TINY: ~2 KB stack that grows and shrinks by copying (lesson 18).
	// 100,000 OS threads would need ~100+ GB of stack; 100,000 goroutines need ~200 MB at
	// most, so launching that many is fine.
	//
	// DATA RACE WARNING: `count++` on a plain int from many goroutines is NOT safe. It is
	// really three steps (read, add, write), and two goroutines can interleave them
	// and lose an update. atomic.Int64 performs the increment as one indivisible CPU
	// operation, so no update is lost (lesson 17 covers atomics and mutexes).
	var count atomic.Int64
	var wg2 sync.WaitGroup
	start := time.Now()
	for range 100_000 { // `100_000`: underscores are allowed in numbers for readability
		wg2.Add(1)
		go func() { // an anonymous function started as a goroutine
			defer wg2.Done()
			count.Add(1)
		}()
	}
	wg2.Wait()
	fmt.Println("goroutines ran:", count.Load(), "elapsed < 2s:", time.Since(start) < 2*time.Second)

	// CLOSURES IN GOROUTINES: each goroutine uses `i` from the loop. Before Go 1.22 all of
	// them shared ONE variable, which had usually reached its final value by the time
	// they ran. Since 1.22 each iteration has its own `i` (lesson 3), so this is correct.
	// Writing to different slice indices from different goroutines is also safe, because
	// no two goroutines touch the same memory. We wait with wg3 before reading results.
	var wg3 sync.WaitGroup
	results := make([]int, 3)
	for i := range 3 {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			results[i] = i * i // each goroutine writes its own slot: no race
		}()
	}
	wg3.Wait()
	fmt.Println(results) // [0 1 4]

	// NumGoroutine: how many exist right now (just main, since the others have finished).
	// GOMAXPROCS(0): how many OS threads may run Go code simultaneously. 0 means "just
	// report the current setting" without changing it. Defaults to the number of CPUs.
	fmt.Println("NumGoroutine:", runtime.NumGoroutine(), "GOMAXPROCS:", runtime.GOMAXPROCS(0))
}

// RULES OF THUMB
//   - Always know how a goroutine will END. One that blocks forever (waiting on a channel
//     nobody will ever send to) is a LEAK: it holds memory until the program exits.
//   - When main returns, the program exits and kills all goroutines, finished or not.
//   - Test with `go run -race ./15-goroutines` to detect data races. The race
//     detector instruments memory accesses and reports unsynchronised ones.
//
// CONCURRENCY vs PARALLELISM
//   Concurrency is STRUCTURING a program as independent tasks that can make progress
//   in overlapping time periods. Parallelism is actually EXECUTING several tasks at
//   the same instant on multiple cores. Goroutines give you concurrency; the scheduler
//   (lesson 19) turns it into parallelism when cores are available.
