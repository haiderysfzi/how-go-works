// Lesson 16: Channels: "share memory by communicating".
package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(n int) <-chan int { // <-chan = receive-only
	out := make(chan int)
	go func() {
		defer close(out) // closing signals "no more values"
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out
}

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
	// Unbuffered: send and receive meet (synchronise).
	ch := make(chan string)
	go func() { ch <- "ping" }()
	fmt.Println(<-ch)

	// Buffered: sends only block when the buffer is full.
	buf := make(chan int, 2)
	buf <- 1
	buf <- 2
	fmt.Println(len(buf), cap(buf))

	// Pipeline: producer -> square -> consumer.
	for v := range square(producer(4)) {
		fmt.Print(v, " ")
	}
	fmt.Println()

	// Receiving from a closed channel gives the zero value and ok=false.
	c := make(chan int, 1)
	c <- 7
	close(c)
	v, ok := <-c
	v2, ok2 := <-c
	fmt.Println(v, ok, v2, ok2)

	// select waits on multiple channel operations.
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

	// Timeout and default.
	select {
	case <-make(chan int):
	case <-time.After(20 * time.Millisecond):
		fmt.Println("timed out")
	}
	select {
	case x := <-make(chan int):
		fmt.Println(x)
	default:
		fmt.Println("nothing ready, did not block")
	}

	// Worker pool: N workers read jobs from one channel.
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	var wg sync.WaitGroup
	for w := 0; w < 3; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				results <- j * 2
			}
		}()
	}
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)
	wg.Wait()
	close(results)
	sum := 0
	for r := range results {
		sum += r
	}
	fmt.Println("pool sum:", sum)

	// A nil channel blocks forever (useful to disable a select case).
	// Sending on a closed channel PANICS. Only the sender should close.
	// A deadlock (all goroutines asleep) is detected and crashes with a clear message.
}
