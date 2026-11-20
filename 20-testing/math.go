// Package-style code that the tests below exercise.
// Run: go test -v ./20-testing   (also: -race, -cover, -bench=.)
package main

import "fmt"

func Add(a, b int) int { return a + b }

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func main() { fmt.Println("see math_test.go, run: go test -v ./20-testing") }
