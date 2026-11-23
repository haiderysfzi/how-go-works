// Lesson 18: Stack vs heap, escape analysis and the garbage collector.
// See the compiler's decisions with: go build -gcflags='-m' ./18-memory-and-gc
//
// WHAT YOU WILL LEARN
//   - The two places a program keeps data: the STACK and the HEAP
//   - How the compiler decides which one (escape analysis)
//   - What the garbage collector does, and how to measure and tune it
//
// BACKGROUND: THE STACK
//
//	Every goroutine has a stack: a block of memory used like a pile of plates. Calling a
//	function pushes a FRAME on top (holding its parameters and local variables);
//	returning pops it. Allocating is nearly free (just move a pointer) and cleanup is
//	free (the frame is simply discarded). The catch: a frame's data must not be
//	needed after the function returns.
//
// BACKGROUND: THE HEAP
//
//	The heap is a large shared region for data that must OUTLIVE the function that
//	created it, or whose size is not known at compile time. Allocating takes more work
//	(find a free spot, record it), and the memory is not reclaimed on return: the
//	GARBAGE COLLECTOR (GC) must find and free it later. Heap use is the main
//	source of GC cost.
//
// In Go you never choose: there is no malloc/free and no `stack`/`heap` keyword.
// The compiler decides for each variable.
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

type Point struct{ X, Y int }

// ESCAPE ANALYSIS, case 1: `p` is used only inside this function, so it stays on the
// STACK. It is released for free when the function returns, and the garbage
// collector never sees it.
func sumStack() int {
	p := Point{1, 2}
	return p.X + p.Y
}

// ESCAPE ANALYSIS, case 2: we return the ADDRESS of the Point. The caller will use it
// after this frame is gone, so the Point "escapes" and the compiler allocates it on
// the HEAP. (`//go:noinline` stops the compiler from pasting this function's body into
// its caller, which could remove the allocation and spoil the measurement below.)
//
//go:noinline
func newPoint() *Point {
	return &Point{3, 4}
}

// runtime.MemStats is a snapshot of memory counters maintained by the runtime:
// HeapAlloc (bytes of live heap), Mallocs (total allocations so far), NumGC (GC cycles).
func memStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

func main() {
	fmt.Println(sumStack(), *newPoint()) // 3 {3 4}

	// MEASURING ALLOCATIONS: `testing` (defined below; unrelated to the standard testing
	// package) runs a function 1000 times and reports how many heap allocations each
	// call made on average. Expect 0 for the stack version and 1 for the heap version.
	// `sink = ...` stores the pointer in a package-level variable so the compiler cannot
	// optimise the call away.
	stackAllocs := testing(func() { _ = sumStack() })
	heapAllocs := testing(func() { sink = newPoint() })
	fmt.Println("allocs per call, stack:", stackAllocs, "heap:", heapAllocs)

	// THE GARBAGE COLLECTOR IN ACTION
	// Step 1: allocate 100 slices of 1 MiB (1<<20 bytes) each = about 100 MiB on the
	// heap, all reachable through `data`.
	before := memStats()
	data := make([][]byte, 0)
	for range 100 {
		data = append(data, make([]byte, 1<<20)) // 100 MiB
	}
	during := memStats()
	// Step 2: `data = nil` drops our only reference. The GC decides what to free by
	// REACHABILITY: anything reachable by following pointers from "roots" (globals,
	// goroutine stacks) is live; everything else is garbage. After this line the 100
	// MiB is garbage, but not yet freed.
	data = nil // drop references: now unreachable
	// Step 3: runtime.GC() forces a collection now (normally it runs automatically).
	// You almost never call this in real code; it is here so the demo is deterministic.
	runtime.GC()
	after := memStats()
	// >>20 divides by 2^20, converting bytes to MiB. Expect "during" far above the others.
	fmt.Printf("heap MiB: before=%d during=%d after=%d\n",
		before.HeapAlloc>>20, during.HeapAlloc>>20, after.HeapAlloc>>20)
	fmt.Println("GC cycles ran:", after.NumGC-before.NumGC >= 1)

	// TUNING
	// GOGC (default 100) is a percentage: start the next collection when the heap has
	// grown 100% past the live data left after the previous collection, i.e. when it has
	// DOUBLED. Higher = fewer, bigger collections (less CPU, more memory); lower = the reverse.
	// SetGCPercent returns the previous setting.
	// SetMemoryLimit sets a SOFT limit on total Go memory (2<<30 = 2 GiB): as the program
	// approaches it, the GC works harder to stay under. Useful in containers with a hard limit.
	old := debug.SetGCPercent(200)
	fmt.Println("previous GOGC:", old)
	debug.SetMemoryLimit(2 << 30)
}

var sink *Point

// testing measures average heap allocations per call of f. Mallocs only counts up, so
// we subtract the "before" count and divide by the number of runs. Rounding hides
// background allocations from the runtime itself.
func testing(f func()) float64 {
	var total float64
	const n = 1000
	before := memStats().Mallocs
	for range n {
		f()
	}
	total = float64(memStats().Mallocs-before) / n
	return float64(int(total + 0.5))
}

// KEY FACTS
//   - HOW THE GC WORKS: it is a "mark-and-sweep" collector. MARK: starting from the
//     roots, follow every pointer and flag everything reachable. SWEEP: reclaim all
//     unflagged memory. "Tri-colour" is the bookkeeping: white (not yet seen), grey
//     (seen, children not yet checked), black (seen and fully checked).
//   - CONCURRENT: marking runs alongside your goroutines, so pauses are typically well
//     under a millisecond. "Non-generational" means it does not treat young objects
//     specially: the whole heap is considered each cycle.
//   - Slices, maps, strings, channels and interfaces all hold pointers into the heap.
//   - To reduce GC work, allocate LESS: pre-size slices (lesson 5), reuse buffers with
//     sync.Pool (lesson 17), avoid needless pointers.
//   - Goroutine stacks start small and grow by copying the whole stack to a bigger
//     block. That is why pointers INTO a stack are never allowed to be stored
//     outside it: they would go stale after a move.
