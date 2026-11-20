// Lesson 19: How the Go runtime schedules goroutines (the G-M-P model).
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// G = goroutine (your code + its stack)
// M = machine, an OS thread
// P = processor, a scheduling context with a local run queue; there are GOMAXPROCS of them
//
// An M needs a P to run Gs. Each P has a queue of runnable Gs; idle Ps STEAL work from busy ones.
// When a G blocks in a syscall, its P is handed to another M, so other goroutines keep running.
// When a G blocks on a channel/mutex/timer it is parked cheaply (no thread blocked).
// Since Go 1.14 the scheduler also preempts long-running goroutines asynchronously.

var sink uint64

// spin does a FIXED amount of CPU work (not "wait until a deadline"), so total time
// really depends on how many cores the scheduler can use.
func spin(iters int) {
	var x uint64 = 1
	for i := 0; i < iters; i++ {
		x = x*6364136223846793005 + 1442695040888963407
	}
	sink = x
}

// calibrate finds an iteration count that takes about 100ms on one core.
func calibrate() int {
	n := 1_000_000
	for {
		start := time.Now()
		spin(n)
		if el := time.Since(start); el > 50*time.Millisecond {
			return int(float64(n) * float64(100*time.Millisecond) / float64(el))
		}
		n *= 2
	}
}

func main() {
	iters := calibrate()
	fmt.Println("cores:", runtime.NumCPU(), "GOMAXPROCS:", runtime.GOMAXPROCS(0))

	run := func(procs int) time.Duration {
		runtime.GOMAXPROCS(procs)
		start := time.Now()
		var wg sync.WaitGroup
		for range 4 {
			wg.Add(1)
			go func() { defer wg.Done(); spin(iters) }()
		}
		wg.Wait()
		return time.Since(start)
	}

	one := run(1)
	many := run(4)
	fmt.Println("4 CPU-bound tasks on 1 P  ~", one.Round(100*time.Millisecond))  // about 400ms
	fmt.Println("4 CPU-bound tasks on 4 Ps ~", many.Round(100*time.Millisecond)) // about 100ms on 4+ cores

	// Blocking tasks do NOT need many Ps: sleeping goroutines are parked.
	runtime.GOMAXPROCS(1)
	start := time.Now()
	var wg sync.WaitGroup
	for range 1000 {
		wg.Add(1)
		go func() { defer wg.Done(); time.Sleep(50 * time.Millisecond) }()
	}
	wg.Wait()
	fmt.Println("1000 sleepers on 1 P ~", time.Since(start).Round(50*time.Millisecond))

	// Cooperative hint (rarely needed now).
	runtime.Gosched()
}
