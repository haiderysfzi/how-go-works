// Lesson 5: Arrays vs slices, and how slices work under the hood.
//
// WHAT YOU WILL LEARN
//   - Arrays: fixed size, copied by value
//   - Slices: a small "window" onto an array, with a length and a capacity
//   - Why slicing and `append` can surprise you (shared memory!)
//   - How to copy, delete, pre-size and build 2D slices
//
// BACKGROUND: memory is one long row of numbered byte slots. An ARRAY is a block of
// same-typed elements stored side by side in one stretch of that row, so element i is
// found by simple arithmetic (start + i*elementSize). That is why indexing is fast.
package main

import "fmt"

func main() {
	// ARRAYS
	// [3]int means "exactly 3 ints". The length is part of the TYPE, so [3]int and
	// [4]int are different, incompatible types. An array is a VALUE: assigning it to
	// another variable copies every element.
	arr := [3]int{1, 2, 3}
	copyArr := arr // a full, independent copy
	copyArr[0] = 99
	fmt.Println(arr, copyArr) // [1 2 3] [99 2 3]  -> arr was not affected
	// Fixed length and full copying make arrays rarely used directly; slices are the norm.

	// SLICES
	// A slice is NOT the data itself. It is a tiny 3-field struct (24 bytes on 64-bit):
	//
	//     +---------+--------+----------+
	//     | pointer | length | capacity |
	//     +----|----+--------+----------+
	//          v
	//     [ 10 | 20 | 30 | 40 | 50 ]      <- a hidden "backing array"
	//
	//   pointer  = address of the first element the slice can see
	//   length   = how many elements you can currently use (what len() returns)
	//   capacity = how many elements fit from `pointer` to the END of the backing array
	//              before more space is needed (what cap() returns)
	s := []int{10, 20, 30, 40, 50}
	fmt.Println(len(s), cap(s)) // 5 5

	// SLICING: s[low:high] makes a NEW slice header pointing INTO THE SAME backing
	// array. Nothing is copied. It covers indices low up to (not including) high.
	//   part: pointer -> s[1], length 2 (elements 20, 30), capacity 4 (from 20 to the end)
	// So writing through `part` changes `s` too, because both look at the same memory.
	part := s[1:3]
	part[0] = 999
	fmt.Println(s, part, len(part), cap(part)) // s changed too; cap(part) is 4

	// APPEND AND GROWTH
	// `append(slice, value)` adds to the end:
	//   - If there is spare capacity, it writes into the existing backing array and
	//     returns a slice with a larger length. Fast, no allocation.
	//   - If the array is FULL, Go allocates a bigger array (roughly double for small
	//     slices), copies the old elements over, adds the new one, and returns a
	//     slice pointing at the NEW array. The old array is left for the garbage collector.
	// You must keep the return value: `grow = append(grow, i)`.
	// The loop prints only when the capacity changes, showing the growth steps.
	var grow []int // a nil slice: length 0, capacity 0
	prevCap := -1
	for i := 0; i < 20; i++ {
		grow = append(grow, i)
		if cap(grow) != prevCap {
			fmt.Print("cap:", cap(grow), " ")
			prevCap = cap(grow)
		}
	}
	fmt.Println()

	// THE CLASSIC GOTCHA
	// make([]int, 3, 10) creates length 3, capacity 10: seven free slots remain.
	// Both appends below find room, so BOTH write into slot index 3 of the SAME array.
	// The second append overwrites the first one's value, and b and c still share it.
	a := make([]int, 3, 10) // len 3, cap 10
	b := append(a, 1)
	c := append(a, 2)       // overwrites the same slot b used!
	fmt.Println(b[3], c[3]) // 2 2

	// FULL SLICE EXPRESSION a[low:high:max] also limits CAPACITY to max-low.
	// With d := a[0:3:3], capacity equals length, so the next append has no room and
	// MUST allocate a new array. That makes `e` independent of `a`.
	d := a[0:3:3]
	e := append(d, 7)
	fmt.Println(len(e), cap(d)) // 4 3

	// PRE-SIZING
	// make([]T, length, capacity). When you know roughly how many elements are coming,
	// reserve the capacity up front to avoid repeated grow-and-copy cycles.
	nums := make([]int, 0, 100) // length 0 (empty), room for 100
	for i := range 5 {
		nums = append(nums, i*i)
	}
	fmt.Println(nums)

	// COPY
	// copy(dst, src) copies min(len(dst), len(src)) elements and returns that count.
	// Unlike slicing, the result lives in separate memory.
	dst := make([]int, 2)
	n := copy(dst, s) // only 2 fit, so only 2 are copied
	fmt.Println(n, dst)

	// DELETING AN ELEMENT: there is no built-in delete. Append the tail onto the head.
	// x[:1] = [1], x[2:] = [3 4], and `...` unpacks a slice into separate arguments.
	// The elements shift left in the SAME backing array (x is modified in place).
	x := []int{1, 2, 3, 4}
	x = append(x[:1], x[2:]...) // delete index 1
	fmt.Println(x)              // [1 3 4]

	// 2D SLICES: a slice whose elements are slices. Each row must be created separately,
	// so rows can even have different lengths ("jagged").
	grid := make([][]int, 2) // 2 rows, each currently a nil slice
	for i := range grid {
		grid[i] = make([]int, 3) // give each row 3 columns
	}
	fmt.Println(grid) // [[0 0 0] [0 0 0]]

	// NIL SLICE vs EMPTY SLICE
	// `var ns []int` has pointer nil; `[]int{}` has a real (zero-length) array.
	// Both have len 0 and both work with append and range, so most code treats them
	// the same. They differ only when compared to nil (and in JSON encoding, lesson 21).
	var ns []int
	es := []int{}
	fmt.Println(ns == nil, es == nil, len(ns), len(es)) // true false 0 0

	// MEMORY LEAK TIP
	// A tiny sub-slice keeps the WHOLE backing array alive, because the garbage
	// collector sees the pointer into it. If you keep 3 elements from a 1 GB array,
	// the 1 GB stays in memory. Copy what you need into a new slice so the big one
	// can be freed.
}
