// Lesson 13: defer, panic, recover.
//
// WHAT YOU WILL LEARN
//   - How `defer` guarantees cleanup, and the exact order deferred calls run
//   - What a panic is, what happens during one ("unwinding"), and what crashes the program
//   - How `recover` stops a panic, and why it only works inside a deferred function
//
// BACKGROUND: the CALL STACK is the chain of functions currently running: main called
// A, A called B, B is running now. A PANIC is Go's emergency stop for situations that
// should never happen (a bug): index out of range, nil pointer use, integer divide by
// zero, or an explicit panic(...) call. When a panic starts, Go UNWINDS the stack:
// the current function stops immediately, its deferred calls run, then the caller
// stops and ITS deferred calls run, and so on up the chain. If nothing recovers,
// the program prints the panic message plus a stack trace and exits with status 2.
package main

import (
	"fmt"
	"os"
)

// NAMED results (result, err) let the deferred function change what the caller sees.
// Flow for safeDivide(1, 0):
//  1. `a / b` panics (integer division by zero).
//  2. Unwinding begins; our deferred closure runs.
//  3. recover() returns the panic value and STOPS the panic. Normal execution resumes
//     at the point where this function returns to its caller.
//  4. The closure stores an error into `err`, which is what the caller receives.
//
// recover() only works when called DIRECTLY inside a deferred function. Anywhere else
// it returns nil and does nothing. If there is no panic, recover() also returns nil,
// which is why we test `r != nil`.
func safeDivide(a, b int) (result int, err error) {
	defer func() {
		// recover only works inside a deferred function.
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	return a / b, nil // integer divide by zero panics
}

// Each pass registers a deferred call. The arguments ("defer", i) are evaluated at the
// moment of the `defer` statement, so each call remembers its own i (0, 1, 2). The calls
// run when loopDefer returns, newest first (LIFO), printing: defer2 defer1 defer0.
// (Deferring inside a loop of thousands of items holds all of them until the function
// ends; for per-item cleanup, call a helper function instead.)
func loopDefer() {
	for i := 0; i < 3; i++ {
		defer fmt.Print("defer", i, " ") // arguments are evaluated NOW, run later (LIFO)
	}
}

// The standard cleanup pattern: acquire a resource, immediately `defer` its release.
// Placing the cleanup next to the acquisition means you cannot forget it, even if a
// later line returns early or panics. LIFO order is what you want here: f.Close()
// runs first, THEN the file is removed.
func cleanup() {
	f, err := os.CreateTemp("", "demo") // creates a real temporary file
	if err != nil {
		return
	}
	defer os.Remove(f.Name()) // runs even if later code panics
	defer f.Close()
	fmt.Fprintln(f, "data") // write a line into the file
	fmt.Println("wrote", f.Name() != "")
}

func main() {
	fmt.Println(safeDivide(10, 2)) // 5 <nil>
	fmt.Println(safeDivide(1, 0))  // 0 recovered: runtime error: integer divide by zero

	loopDefer()
	fmt.Println()

	cleanup()

	// DEFERRED CALLS RUN DURING A PANIC, before the stack unwinds to the next function
	// up. Inside this anonymous function (called immediately by the trailing `()`):
	//   - two defers are registered
	//   - panic("something broke") begins unwinding
	//   - defers run LIFO: first the recover closure (prints "recovered: something broke"
	//     and ends the panic), then the Println ("deferred runs during panic")
	//   - the function returns normally to main, which carries on
	func() {
		defer fmt.Println("deferred runs during panic")
		defer func() { fmt.Println("recovered:", recover()) }()
		panic("something broke")
	}()

	// Runtime panics: out-of-range index, writing to a nil map, nil pointer dereference.
	// This defer is the safety net for `s[3]` below: s is empty, so index 3 is invalid.
	// It runs as main unwinds, printing the panic message; main then returns normally.
	defer func() { fmt.Println("caught:", recover()) }()
	var s []int
	_ = s[3]
}

// GUIDANCE
//   - Return errors for EXPECTED failures (file missing, bad input, network down).
//   - Panic only for impossible states and programmer mistakes.
//   - Libraries should not let panics escape their public API: recover at the boundary
//     and convert to an error (like safeDivide does).
//   - A panic inside a goroutine (lesson 15) that is not recovered INSIDE THAT
//     goroutine crashes the WHOLE program; recover in main does not catch it.
