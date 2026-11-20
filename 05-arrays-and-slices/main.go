// Lesson 5: Arrays vs slices, and how slices work under the hood.
package main

import "fmt"

func main() {
	// Arrays have a FIXED length that is part of the type, and are copied by value.
	arr := [3]int{1, 2, 3}
	copyArr := arr
	copyArr[0] = 99
	fmt.Println(arr, copyArr) // [1 2 3] [99 2 3]

	// A slice is a small struct: {pointer to array, length, capacity}.
	s := []int{10, 20, 30, 40, 50}
	fmt.Println(len(s), cap(s))

	// Slicing shares the same backing array: NO copy.
	part := s[1:3]
	part[0] = 999
	fmt.Println(s, part, len(part), cap(part)) // s changed too; cap(part) is 4

	// append may reallocate when capacity runs out.
	var grow []int
	prevCap := -1
	for i := 0; i < 20; i++ {
		grow = append(grow, i)
		if cap(grow) != prevCap {
			fmt.Print("cap:", cap(grow), " ")
			prevCap = cap(grow)
		}
	}
	fmt.Println()

	// The classic gotcha: append may or may not share memory with the original.
	a := make([]int, 3, 10) // len 3, cap 10
	b := append(a, 1)
	c := append(a, 2)       // overwrites the same slot b used!
	fmt.Println(b[3], c[3]) // 2 2

	// Full slice expression limits capacity to force a copy on append.
	d := a[0:3:3]
	e := append(d, 7)
	fmt.Println(len(e), cap(d))

	// Pre-size when you know the length, to avoid reallocations.
	nums := make([]int, 0, 100)
	for i := range 5 {
		nums = append(nums, i*i)
	}
	fmt.Println(nums)

	// copy copies min(len(dst), len(src)) elements into independent memory.
	dst := make([]int, 2)
	n := copy(dst, s)
	fmt.Println(n, dst)

	// Delete element i, insert, and 2D slices.
	x := []int{1, 2, 3, 4}
	x = append(x[:1], x[2:]...) // delete index 1
	fmt.Println(x)
	grid := make([][]int, 2)
	for i := range grid {
		grid[i] = make([]int, 3)
	}
	fmt.Println(grid)

	// nil slice vs empty slice: both have len 0 and both work with append/range.
	var ns []int
	es := []int{}
	fmt.Println(ns == nil, es == nil, len(ns), len(es))

	// Memory leak tip: a tiny sub-slice keeps the WHOLE backing array alive.
	// Copy out what you need if the original is large.
}
