// Lesson 13: defer, panic, recover.
package main

import (
	"fmt"
	"os"
)

func safeDivide(a, b int) (result int, err error) {
	defer func() {
		// recover only works inside a deferred function.
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	return a / b, nil // integer divide by zero panics
}

func loopDefer() {
	for i := 0; i < 3; i++ {
		defer fmt.Print("defer", i, " ") // arguments are evaluated NOW, run later (LIFO)
	}
}

func cleanup() {
	f, err := os.CreateTemp("", "demo")
	if err != nil {
		return
	}
	defer os.Remove(f.Name()) // runs even if later code panics
	defer f.Close()
	fmt.Fprintln(f, "data")
	fmt.Println("wrote", f.Name() != "")
}

func main() {
	fmt.Println(safeDivide(10, 2))
	fmt.Println(safeDivide(1, 0))

	loopDefer()
	fmt.Println()

	cleanup()

	// Order matters: deferred calls run BEFORE the panic unwinds the next frame up.
	func() {
		defer fmt.Println("deferred runs during panic")
		defer func() { fmt.Println("recovered:", recover()) }()
		panic("something broke")
	}()

	// Out-of-range index, nil map writes and nil derefs also panic.
	defer func() { fmt.Println("caught:", recover()) }()
	var s []int
	_ = s[3]
}

// Guidance: return errors for expected failures. Panic for impossible states.
// Libraries should not let panics escape their API. A panic in a goroutine you did
// not recover crashes the WHOLE program.
