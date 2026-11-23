// Lesson 16: Channels: "share memory by communicating".
//
// WHAT YOU WILL LEARN
//   - What a channel is, and the difference between unbuffered and buffered ones
//   - Send, receive, close, and `range` over a channel
//   - Directional channel types (<-chan, chan<-) and pipelines
//   - `select` for waiting on several things, with timeouts and non-blocking tries
//   - The worker-pool pattern
//
// BACKGROUND: goroutines (lesson 15) need a safe way to hand data to each other.
// Go's motto: "Don't communicate by sharing memory; share memory by communicating."
// A CHANNEL is a typed pipe between goroutines: one side SENDS a value (`ch <- v`),
// another RECEIVES it (`v := <-ch`). The runtime handles all locking internally, so
// only one goroutine "owns" a value at a time and data races are avoided.
package main

import (
	"fmt"
	"sync"
	"time"
)

// producer returns a RECEIVE-ONLY channel: `<-chan int` means callers can read from it
// but not send to or close it. The goroutine inside holds the full `chan int` (so it
// can send and close). Restricting direction in a signature documents who does what
// and the compiler enforces it.
// `defer close(out)` runs when the loop is finished; closing tells receivers "no more
// values are coming".
func producer(n int) <-chan int { // <-chan = receive-only
	out := make(chan int) // UNBUFFERED: each send waits for a matching receive
	go func() {
		defer close(out) // closing signals "no more values"
		for i := 1; i <= n; i++ {
			out <- i // blocks until the consumer takes the value
		}
	}()
	return out
}

// A PIPELINE STAGE: reads from one channel, does work, writes to another. Chaining
// stages (producer -> square -> consumer) lets each run concurrently.
// `for v := range in` receives repeatedly and exits automatically once `in` is
// closed and drained. We then close OUR output so the next stage stops too.
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in { // range stops when the channel is closed
			out <- v * v
		}
	}()
	return out
}

func main() {
	// UNBUFFERED CHANNEL (capacity 0): a send does not complete until a receiver is
	// ready, and vice versa: the two goroutines "meet" at that moment, which also
	// SYNCHRONISES them. Here a goroutine sends "ping" and main waits on <-ch.
	ch := make(chan string)
	go func() { ch <- "ping" }()
	fmt.Println(<-ch)

	// BUFFERED CHANNEL: make(chan T, N) has room for N values. A send blocks only when
	// the buffer is FULL; a receive blocks only when it is EMPTY. Think of N as a small
	// queue between sender and receiver. len = items waiting, cap = buffer size.
	buf := make(chan int, 2)
	buf <- 1
	buf <- 2                        // a third send here would block forever (deadlock) since nobody receives
	fmt.Println(len(buf), cap(buf)) // 2 2

	// Pipeline: producer sends 1..4, square squares each, main prints them.
	for v := range square(producer(4)) {
		fmt.Print(v, " ") // 1 4 9 16
	}
	fmt.Println()

	// RECEIVING FROM A CLOSED CHANNEL: values still in the buffer are delivered first.
	// After it is empty, every receive returns immediately with the type's ZERO VALUE
	// and ok=false. The "comma ok" lets you distinguish a real 0 from "closed".
	c := make(chan int, 1)
	c <- 7
	close(c)
	v, ok := <-c   // 7 true
	v2, ok2 := <-c // 0 false
	fmt.Println(v, ok, v2, ok2)

	// SELECT: like a switch, but each case is a channel operation. It waits until at
	// least one case can proceed, then runs that one. If several are ready, it picks
	// one at random (fairness). Here "a" is ready after 10 ms and "b" after 30 ms, so
	// the loop (which runs twice) prints a's message first and b's second.
	a, b := make(chan string), make(chan string)
	go func() { time.Sleep(10 * time.Millisecond); a <- "from a" }()
	go func() { time.Sleep(30 * time.Millisecond); b <- "from b" }()
	for range 2 {
		select {
		case m := <-a:
			fmt.Println(m)
		case m := <-b:
			fmt.Println(m)
		}
	}

	// TIMEOUT: time.After(d) returns a channel that receives one value after d elapses.
	// The first case waits on a channel nobody ever sends to, so only the timer fires.
	// This is the standard way to stop waiting on something slow.
	select {
	case <-make(chan int):
	case <-time.After(20 * time.Millisecond):
		fmt.Println("timed out")
	}
	// DEFAULT: a `default` case runs immediately if no other case is ready, making
	// the operation NON-BLOCKING ("try, and move on if not possible").
	select {
	case x := <-make(chan int):
		fmt.Println(x)
	default:
		fmt.Println("nothing ready, did not block")
	}

	// WORKER POOL: a fixed number of workers (3) all read from the same `jobs`
	// channel. Whichever worker is free grabs the next job, so work is shared automatically
	// and concurrency is capped at 3 regardless of how many jobs exist.
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	var wg sync.WaitGroup
	for w := 0; w < 3; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs { // each worker loops until jobs is closed and empty
				results <- j * 2
			}
		}()
	}
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)    // tells the workers: no more jobs, finish up
	wg.Wait()      // all workers have exited, so every result has been sent
	close(results) // now we can close results so the range below terminates
	sum := 0
	for r := range results {
		sum += r
	}
	fmt.Println("pool sum:", sum) // 2+4+6+8+10 = 30

	// FACTS WORTH MEMORISING
	// - A nil channel blocks forever on send and receive. In a select, that effectively
	//   DISABLES the case (set a channel variable to nil to turn it off).
	// - Sending on a closed channel PANICS. Closing an already-closed channel panics.
	//   Only the sender should close, and only when it is sure nobody else will send.
	//   Closing is optional unless a receiver needs to know the stream ended (`range`).
	// - A DEADLOCK (every goroutine blocked, nothing can ever wake them) is detected by
	//   the runtime and the program crashes with "all goroutines are asleep - deadlock!".
}
