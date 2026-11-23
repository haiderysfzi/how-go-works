// Lesson 19: How the Go runtime schedules goroutines (the G-M-P model).
//
// WHAT YOU WILL LEARN
//   - What a scheduler is and why Go has its own
//   - The three actors: G (goroutine), M (OS thread), P (processor/context)
//   - Work stealing, blocking, and preemption
//   - How to SEE the effect of GOMAXPROCS with measurements
//
// BACKGROUND: the operating system already schedules THREADS onto CPU cores. But OS
// threads are expensive to create and switch (kernel involvement, megabytes of
// stack). Go therefore runs its OWN scheduler in user space: it multiplexes a large
// number of cheap goroutines (G) onto a small number of OS threads (M). This is
// called M:N scheduling. A switch between goroutines costs on the order of tens
// of nanoseconds, with no kernel call.
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// THE G-M-P MODEL
//
//	G = Goroutine: your function plus its small stack and bookkeeping (state, what it waits on).
//	M = Machine:   a real OS thread. Only an M actually executes code on a CPU core.
//	P = Processor: a SCHEDULING CONTEXT (a token + resources, not hardware). It holds a LOCAL
//	               RUN QUEUE of goroutines that are ready to run, and a memory cache. There
//	               are exactly GOMAXPROCS of them (default: number of CPU cores).
//
//	  P0 queue: [G4 G5 G6]       P1 queue: [G7]        P2 queue: []   <- idle, will steal
//	     |                          |                     |
//	     M0 (running G3)            M1 (running G8)       M2
//
// RULES
//   - An M must HOLD a P to run Go code. So at most GOMAXPROCS goroutines run in parallel.
//   - When a goroutine finishes or blocks, its M takes the next G from its P's queue.
//   - WORK STEALING: a P with an empty queue steals about half the goroutines from a
//     busy P's queue (or from a global queue), so cores don't sit idle while others
//     have backlog.
//   - SYSCALL BLOCKING: when a G makes a blocking system call (e.g. reading a file),
//     its M is stuck in the kernel, so the runtime detaches the P and hands it to
//     another M to keep running other goroutines. When the call returns, the G
//     queues up for a P again.
//   - PARKING: when a G blocks on a channel, mutex, timer or network I/O, it is
//     PARKED: marked waiting and removed from the run queue. NO thread is blocked and
//     almost no resources are used. It is made runnable again when the event happens.
//   - PREEMPTION: since Go 1.14, a goroutine that runs for more than ~10ms without
//     yielding is interrupted asynchronously (via a signal to its thread), so one tight
//     loop cannot starve the others. Before that, the loop had to call a function
//     to give the scheduler a chance.

// sink is a package-level variable that spin writes its result into. Because the result
// is "used" outside the function, the compiler cannot delete the loop as dead code.
var sink uint64

// spin does a FIXED amount of CPU work (not "wait until a deadline"), so total time
// really depends on how many cores the scheduler can use.
// The loop keeps multiplying and adding large constants (a simple pseudo-random
// generator). Unsigned integer overflow wraps (lesson 2), which is fine here: we
// just need work that cannot be skipped.
func spin(iters int) {
	var x uint64 = 1
	for i := 0; i < iters; i++ {
		x = x*6364136223846793005 + 1442695040888963407
	}
	sink = x
}

// calibrate finds an iteration count that takes about 100ms on one core.
// Method: time spin(n); if it took under 50ms the sample is too small to trust, so
// double n and retry. Once it takes >50ms, scale n proportionally to land on 100ms
// (n * 100ms / elapsed). The result differs per machine, which is why we measure
// rather than hard-code a number.
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
	// NumCPU = logical cores. GOMAXPROCS(0) = how many Ps exist now (0 = only report).
	fmt.Println("cores:", runtime.NumCPU(), "GOMAXPROCS:", runtime.GOMAXPROCS(0))

	// `run` is a closure (lesson 4). It sets the number of Ps, starts 4 goroutines that
	// each burn ~100ms of CPU, waits for all of them, and returns the elapsed wall-clock time.
	run := func(procs int) time.Duration {
		runtime.GOMAXPROCS(procs) // change the number of Ps (this briefly stops the world; fine for a demo)
		start := time.Now()
		var wg sync.WaitGroup
		for range 4 {
			wg.Add(1)
			go func() { defer wg.Done(); spin(iters) }()
		}
		wg.Wait()
		return time.Since(start)
	}

	// EXPERIMENT 1: CPU-bound work.
	// With 1 P only ONE goroutine can execute at a time, so four 100ms tasks take
	// ~400ms (they take turns). With 4 Ps (and 4+ cores) all four run at once: ~100ms.
	// Adding Ps helps only when the work actually uses the CPU, and only up to your core count.
	one := run(1)
	many := run(4)
	fmt.Println("4 CPU-bound tasks on 1 P  ~", one.Round(100*time.Millisecond))  // about 400ms
	fmt.Println("4 CPU-bound tasks on 4 Ps ~", many.Round(100*time.Millisecond)) // about 100ms on 4+ cores

	// EXPERIMENT 2: waiting work.
	// 1000 goroutines each sleep 50ms with ONE P. A sleeping goroutine is PARKED (it
	// uses no P and no thread), so all 1000 wait at the same time and the whole thing
	// still takes ~50ms, not 50 seconds. This is why one Go server can hold tens of
	// thousands of mostly-idle connections.
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

	// runtime.Gosched() voluntarily yields the processor: "let other goroutines run,
	// then put me back in the queue". It was needed in the days before preemption;
	// now it is rarely useful.
	runtime.Gosched()
}
