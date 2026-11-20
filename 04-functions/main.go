// Lesson 4: Functions, multiple returns, closures, variadics, defer.
package main

import (
	"errors"
	"fmt"
)

// Multiple return values are idiomatic; the last is usually an error.
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}

// Named results act as pre-declared variables; a bare return returns them.
func minMax(nums ...int) (min, max int) { // ...int is variadic
	min, max = nums[0], nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return
}

// Functions are values.
func apply(fn func(int) int, v int) int { return fn(v) }

// Closures capture variables by reference.
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// defer runs when the surrounding function returns, in LIFO order.
func deferDemo() (result int) {
	defer fmt.Println("deferred 1")
	defer fmt.Println("deferred 2 (runs first)")
	defer func() { result *= 2 }() // can modify named results
	return 21
}

func main() {
	if q, err := divide(10, 4); err == nil {
		fmt.Println(q)
	}
	_, err := divide(1, 0)
	fmt.Println("error:", err)

	lo, hi := minMax(4, 9, 2, 7)
	fmt.Println(lo, hi)
	nums := []int{5, 1}
	fmt.Println(minMax(nums...)) // spread a slice into a variadic

	fmt.Println(apply(func(x int) int { return x * x }, 6))

	next := counter()
	next()
	next()
	fmt.Println(next()) // 3

	fmt.Println(deferDemo()) // 42

	// Function types and methods on them.
	var h HandlerFunc = func(s string) { fmt.Println("handled", s) }
	h.Serve("x")

	// Anonymous function invoked immediately.
	func() { fmt.Println("IIFE") }()
}

type HandlerFunc func(string)

func (f HandlerFunc) Serve(s string) { f(s) } // this is how net/http's HandlerFunc works

// Notes: Go passes EVERYTHING by value (arguments are copies). Slices, maps and
// pointers copy a small header/pointer, which is why they appear to be "by reference".
// No default arguments, no overloading, no optional params.
