// Lesson 15: Goroutines: cheap, runtime-managed concurrent functions.
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(id) * time.Millisecond)
	fmt.Println("worker", id, "done")
}

func main() {
	// `go f()` starts f concurrently. main does NOT wait for it.
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait() // block until the counter hits zero

	// Goroutines are tiny (~2KB starting stack, growable), so 100k is fine.
	var count atomic.Int64
	var wg2 sync.WaitGroup
	start := time.Now()
	for range 100_000 {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			count.Add(1)
		}()
	}
	wg2.Wait()
	fmt.Println("goroutines ran:", count.Load(), "elapsed < 2s:", time.Since(start) < 2*time.Second)

	// Closures in goroutines: Go 1.22 gives each iteration its own i.
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
	fmt.Println(results)

	fmt.Println("NumGoroutine:", runtime.NumGoroutine(), "GOMAXPROCS:", runtime.GOMAXPROCS(0))
}

// Rules of thumb:
// - Always know how a goroutine will END, otherwise it leaks.
// - When main returns, the program exits and kills all goroutines.
// - Test with `go run -race ./15-goroutines` to detect data races.
// Concurrency is not parallelism: concurrency is structuring work as independent tasks,
// parallelism is running them at the same instant on multiple cores.
