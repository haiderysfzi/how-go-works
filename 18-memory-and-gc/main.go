// Lesson 18: Stack vs heap, escape analysis and the garbage collector.
// See the compiler's decisions with: go build -gcflags='-m' ./18-memory-and-gc
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

type Point struct{ X, Y int }

// Does not escape: p lives on the stack, freed for free when the function returns.
func sumStack() int {
	p := Point{1, 2}
	return p.X + p.Y
}

// Escapes: returning a pointer forces allocation on the heap.
//
//go:noinline
func newPoint() *Point {
	return &Point{3, 4}
}

func memStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

func main() {
	fmt.Println(sumStack(), *newPoint())

	// Measure allocations: the stack version allocates nothing.
	stackAllocs := testing(func() { _ = sumStack() })
	heapAllocs := testing(func() { sink = newPoint() })
	fmt.Println("allocs per call, stack:", stackAllocs, "heap:", heapAllocs)

	// GC in action.
	before := memStats()
	data := make([][]byte, 0)
	for range 100 {
		data = append(data, make([]byte, 1<<20)) // 100 MiB
	}
	during := memStats()
	data = nil // drop references: now unreachable
	runtime.GC()
	after := memStats()
	fmt.Printf("heap MiB: before=%d during=%d after=%d\n",
		before.HeapAlloc>>20, during.HeapAlloc>>20, after.HeapAlloc>>20)
	fmt.Println("GC cycles ran:", after.NumGC-before.NumGC >= 1)

	// Tunables: GOGC (default 100) = collect when heap doubles since last GC.
	// A soft memory limit helps in containers.
	old := debug.SetGCPercent(200)
	fmt.Println("previous GOGC:", old)
	debug.SetMemoryLimit(2 << 30)
}

var sink *Point

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

// Key facts:
// - Go's GC is concurrent, non-generational, tri-colour mark-and-sweep. Pauses are
//   typically well under a millisecond.
// - Slices, maps, strings, channels and interfaces point into the heap.
// - Reduce GC work by allocating less: pre-size slices, reuse buffers (sync.Pool),
//   avoid needless pointers.
// - Goroutine stacks start small and grow by copying, so pointers into a stack are never
//   kept outside it.
