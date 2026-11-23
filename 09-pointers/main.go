// Lesson 9: Pointers. Like C pointers, without pointer arithmetic.
//
// WHAT YOU WILL LEARN
//   - What a memory address is, and what a pointer stores
//   - The two pointer operators: & (address-of) and * (dereference)
//   - Why Go copies arguments, and how pointers let a function change the caller's data
//   - Why returning a pointer to a local variable is safe in Go
//   - What a nil pointer is, and why using one crashes
//
// BACKGROUND: every variable lives somewhere in memory, and every byte of memory has a
// numbered ADDRESS (like a house number on a very long street). A POINTER is a variable
// whose value is the address of another variable. Instead of holding the data, it
// holds where the data is. Picture:
//
//	x := 1          x   [ address 0x1000 | value 1      ]
//	p := &x         p   [ address 0x2000 | value 0x1000 ]  --> points at x
package main

import "fmt"

type Big struct{ data [4]int }

// Go passes arguments BY VALUE: the function receives its own COPY. So `n++` below
// increments the copy, and the caller's variable never changes.
func incrementCopy(n int) { n++ }

// `*int` means "pointer to an int". The function receives a copy of the ADDRESS, but
// both copies point at the same int. `*n++` follows the address and increments the int
// found there (`*n` reads "the value that n points to").
func incrementPtr(n *int) { *n++ }

// For a pointer to a struct, `b.data` is shorthand for `(*b).data`; Go auto-dereferences.
// Passing *Big also avoids copying the whole struct on every call.
func modifyBig(b *Big) { b.data[0] = 1 } // field access auto-dereferences

// ESCAPE ANALYSIS: in C this would be a bug, because `x` lives in this function's stack
// frame, which is destroyed on return, leaving a dangling pointer. In Go the compiler
// notices the address of x outlives the function and puts x on the HEAP instead (a
// region the garbage collector manages; see lesson 18). So this is safe.
func newInt() *int { x := 5; return &x } // safe: x escapes to the heap

func main() {
	x := 1
	incrementCopy(x)
	fmt.Println(x)   // 1: the function got a copy
	incrementPtr(&x) // &x = "the address of x"
	fmt.Println(x)   // 2

	p := &x // p has type *int and holds the address of x
	*p = 10 // dereference to write: store 10 into the variable p points to (x)
	// x and *p are two names for the same memory. `p == &x` compares addresses: true.
	fmt.Println(x, *p, p == &x) // 10 10 true

	// new(T) allocates memory for a T, fills it with the zero value, and returns a
	// pointer to it. `q := new(int)` is the same as `var tmp int; q := &tmp`.
	q := new(int)
	*q = 3
	fmt.Println(*q)

	// Returning a pointer to a local variable is fine (unlike C): Go moves it to the heap.
	fmt.Println(*newInt())

	b := Big{}
	modifyBig(&b)  // pass the address so the function edits OUR b, not a copy
	fmt.Println(b) // {[1 0 0 0]}

	// NIL POINTERS
	// The zero value of a pointer is nil: "points at nothing". Comparing to nil is
	// allowed, but following (dereferencing) a nil pointer is a runtime error: there is
	// no struct at that address to read. Go stops with a PANIC rather than quietly
	// reading random memory. We catch the panic with defer+recover (lesson 13) so
	// the rest of the demo can finish and print the message.
	var np *Big
	fmt.Println(np == nil) // true
	defer func() { fmt.Println("recovered:", recover()) }()
	fmt.Println(np.data) // panic: invalid memory address or nil pointer dereference
}

// WHEN TO USE POINTERS
// - You need to mutate the caller's value.
// - The struct is large and copying is expensive.
// - You need to express "optional" (nil means absent).
// - Consistency: if one method needs a pointer receiver, use pointers for all of them.
//
// No pointer arithmetic: in C you can write p+1 to step through memory, which is
// the source of many security bugs. Go forbids it (except through the `unsafe`
// package, which you should avoid).
