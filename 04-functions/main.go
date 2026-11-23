// Lesson 4: Functions, multiple returns, closures, variadics, defer.
//
// WHAT YOU WILL LEARN
//   - How to declare functions and return several values
//   - Named results and "bare" returns
//   - Variadic functions (any number of arguments)
//   - Functions as values, and closures (functions that remember variables)
//   - `defer`, which schedules cleanup for when a function ends
//   - Methods on function types (how net/http's HandlerFunc works)
//
// BACKGROUND: when you call a function, the runtime reserves a block of memory called a
// STACK FRAME for that call's parameters and local variables (lesson 18 goes deeper).
// When the function returns, the frame is released. Keep that picture in mind below.
package main

import (
	"errors"
	"fmt"
)

// MULTIPLE RETURN VALUES
// A function can return more than one value. Go has no exceptions for ordinary
// failures, so the convention is: return the result AND an `error` as the last value.
// The caller checks `err != nil` (lesson 12). `nil` means "no error".
// Parameter syntax `a, b float64` is shorthand for `a float64, b float64`.
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}

// NAMED RESULTS + VARIADIC PARAMETERS
//   - `nums ...int` is VARIADIC: the caller may pass any number of ints, and inside
//     the function `nums` is simply a slice ([]int) holding them all.
//   - `(min, max int)` NAMES the results. Named results are real variables, declared
//     for you and set to their zero value (0) at the start of the function.
//   - A bare `return` (no values) returns whatever the named variables currently hold.
//     Handy in short functions; confusing in long ones, so use it sparingly.
//
// (`min` and `max` shadow Go's built-in functions of the same names. That is legal
// but you would normally pick other names. It is fine for a demo.)
func minMax(nums ...int) (min, max int) { // ...int is variadic
	min, max = nums[0], nums[0]
	for _, n := range nums[1:] { // nums[1:] = every element except the first
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return
}

// FUNCTIONS ARE VALUES
// A function has a type, written from its parameter and result types: func(int) int.
// You can store functions in variables, pass them as arguments and return them,
// exactly like an int or a string. `apply` takes a function `fn` and calls it.
func apply(fn func(int) int, v int) int { return fn(v) }

// CLOSURES
// The function returned here uses `count`, a variable that belongs to the OUTER
// function. A function that "closes over" outside variables like this is a CLOSURE.
//
// Normally `count` would be destroyed when counter() returns. But the returned
// function still needs it, so Go's compiler moves `count` to the HEAP (lesson 18)
// where it outlives the call. The closure holds a reference to that one variable, not
// a copy, which is why every call increments the SAME counter.
// Each call to counter() creates a brand-new, independent `count`.
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// DEFER
// `defer someCall()` does not run the call now. It schedules it to run right BEFORE the
// surrounding function returns, no matter how it exits (normal return, early return,
// or panic). It is the standard way to guarantee cleanup: closing files, unlocking
// mutexes, and so on.
//   - Deferred calls run in LIFO order (last deferred, first run), like a stack of plates.
//   - Arguments of a deferred call are evaluated IMMEDIATELY, only the call is delayed.
//   - Because `result` is a NAMED result, a deferred closure can still change it after
//     the `return` statement has set it. Order of events here:
//     1. `return 21` stores 21 into `result`
//     2. deferred funcs run: first the closure (result becomes 42), then the two Printlns
//     3. the function hands back `result`, which is now 42
func deferDemo() (result int) {
	defer fmt.Println("deferred 1")
	defer fmt.Println("deferred 2 (runs first)")
	defer func() { result *= 2 }() // can modify named results
	return 21
}

func main() {
	// `if` with an init statement (lesson 3): call divide, then test the error.
	if q, err := divide(10, 4); err == nil {
		fmt.Println(q) // 2.5
	}
	// `_` discards the first result because we only care about the error here.
	_, err := divide(1, 0)
	fmt.Println("error:", err)

	lo, hi := minMax(4, 9, 2, 7) // four separate arguments are packed into one slice
	fmt.Println(lo, hi)
	nums := []int{5, 1}
	// `nums...` does the reverse: it UNPACKS an existing slice into variadic arguments.
	fmt.Println(minMax(nums...))

	// A function literal (also called an anonymous function) is written inline
	// and passed straight to apply.
	fmt.Println(apply(func(x int) int { return x * x }, 6)) // 36

	next := counter()   // `next` is a closure with its own private `count`
	next()              // 1
	next()              // 2
	fmt.Println(next()) // 3

	fmt.Println(deferDemo()) // 42

	// FUNCTION TYPES CAN HAVE METHODS
	// `HandlerFunc` (defined below) is a named type whose underlying type is a function.
	// We assign a function literal to it and then call a METHOD on that function value.
	var h HandlerFunc = func(s string) { fmt.Println("handled", s) }
	h.Serve("x")

	// An anonymous function that is defined and called in the same breath: the
	// trailing () is the call. Useful for scoping a few lines or with `defer`/`go`.
	func() { fmt.Println("IIFE") }()
}

// HandlerFunc is a named type for "a function taking a string and returning nothing".
type HandlerFunc func(string)

// Serve is a method on that function type; it simply calls the function itself.
// This trick lets an ordinary function satisfy an interface that requires a method.
// It is exactly how net/http's `http.HandlerFunc` lets you pass a plain function
// where an `http.Handler` (an interface with a ServeHTTP method) is expected.
func (f HandlerFunc) Serve(s string) { f(s) } // this is how net/http's HandlerFunc works

// NOTES
//   - Go passes EVERYTHING by value: when you call f(x), f receives a COPY of x.
//     Slices, maps and pointers are small headers/addresses, so the copy still
//     points at the same underlying data. That is why they look "by reference";
//     the rule is still pass-by-value (lessons 5, 6 and 9).
//   - There are no default argument values, no function overloading (two functions
//     with the same name and different parameters), and no optional parameters.
//     Go chooses simplicity: one name, one meaning.
