// Package-style code that the tests in math_test.go exercise.
// Run: go test -v ./20-testing   (also: -race, -cover, -bench=.)
//
// HOW GO TESTING WORKS
//
//	Testing is built into the toolchain: there is no extra framework to install.
//	`go test` looks for files whose names end in `_test.go`, compiles them together with
//	the package, and runs every function that matches a naming pattern:
//	  TestXxx(t *testing.T)       -> a normal test
//	  BenchmarkXxx(b *testing.B)  -> a performance measurement
//	  FuzzXxx(f *testing.F)       -> a fuzz test (random input generator)
//	Files ending in _test.go are NOT included in a normal `go build`, so tests
//	never bloat your program.
//
// THIS FILE holds the ordinary functions under test. They are intentionally tiny so the
// focus stays on how tests are written; see math_test.go.
package main

import "fmt"

// Add returns the sum of two ints.
func Add(a, b int) int { return a + b }

// Abs returns the absolute value (distance from zero) of n.
// Edge case worth knowing: for the smallest possible int (math.MinInt), -n overflows
// and gives the same negative number back. The fuzz test in math_test.go mentions this.
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// A package main needs a main function to build; this one just points you to the tests.
func main() { fmt.Println("see math_test.go, run: go test -v ./20-testing") }
