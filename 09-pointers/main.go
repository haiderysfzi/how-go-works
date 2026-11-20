// Lesson 9: Pointers. Like C pointers, without pointer arithmetic.
package main

import "fmt"

type Big struct{ data [4]int }

func incrementCopy(n int) { n++ }
func incrementPtr(n *int) { *n++ }
func modifyBig(b *Big)    { b.data[0] = 1 }     // field access auto-dereferences
func newInt() *int        { x := 5; return &x } // safe: x escapes to the heap

func main() {
	x := 1
	incrementCopy(x)
	fmt.Println(x) // 1: the function got a copy
	incrementPtr(&x)
	fmt.Println(x) // 2

	p := &x // p holds the address of x
	*p = 10 // dereference to write
	fmt.Println(x, *p, p == &x)

	// new allocates a zeroed value and returns a pointer.
	q := new(int)
	*q = 3
	fmt.Println(*q)

	// Returning a pointer to a local variable is fine (unlike C): Go moves it to the heap.
	fmt.Println(*newInt())

	b := Big{}
	modifyBig(&b)
	fmt.Println(b)

	// nil pointer dereference panics.
	var np *Big
	fmt.Println(np == nil)
	defer func() { fmt.Println("recovered:", recover()) }()
	fmt.Println(np.data) // panic: invalid memory address or nil pointer dereference
}

// When to use pointers:
// - You need to mutate the caller's value.
// - The struct is large and copying is expensive.
// - You need to express "optional" (nil means absent).
// - Consistency: if one method needs a pointer receiver, use pointers for all of them.
// No pointer arithmetic (except via the unsafe package, which you should avoid).
