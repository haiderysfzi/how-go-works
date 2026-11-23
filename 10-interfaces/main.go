// Lesson 10: Interfaces are satisfied IMPLICITLY (structural typing).
//
// WHAT YOU WILL LEARN
//   - What an interface is and why it is useful
//   - Implicit satisfaction (no `implements` keyword)
//   - Type assertions, the empty interface, and method sets
//   - The famous "nil interface" gotcha, and how interfaces look in memory
//
// BACKGROUND: an INTERFACE is a contract written as a list of method signatures:
// "anything that can do these things." Code written against the interface works with
// every type that satisfies it, so you can swap implementations (a real database vs a
// fake in tests, a file vs a network connection) without changing the caller.
package main

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// An interface is a set of method signatures. No "implements" keyword.
// In Java you must declare `class Circle implements Shape`. In Go, ANY type that has
// the methods Area() float64 and Perimeter() float64 automatically IS a Shape. The
// compiler checks this at the place you use it. This is called structural typing
// (also "duck typing", but verified at compile time). It lets you define interfaces
// AFTER the types, even for types in other people's packages.
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct{ R float64 }
type Square struct{ S float64 }

func (c Circle) Area() float64      { return math.Pi * c.R * c.R }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.R }
func (s Square) Area() float64      { return s.S * s.S }
func (s Square) Perimeter() float64 { return 4 * s.S }

// "Accept interfaces, return concrete types": make parameters as general as possible
// so callers have freedom, but return real types so callers get every feature.
// `describe` works with ANY Shape. %T prints the concrete type hiding in the interface.
func describe(s Shape) string {
	return fmt.Sprintf("%T area=%.2f perimeter=%.2f", s, s.Area(), s.Perimeter())
}

// Small interfaces are best. io.Reader has ONE method, `Read(p []byte) (n int, err error)`,
// meaning "fill this buffer with some bytes". Files, network connections, HTTP bodies,
// compressed streams, strings and more are all Readers, so one function like this can
// read from any of them. io.Copy pulls bytes from the Reader until it ends; io.Discard
// is a writer that throws everything away. Net effect: it counts the bytes.
func countBytes(r io.Reader) int {
	n, _ := io.Copy(io.Discard, r)
	return int(n)
}

type Animal interface{ Sound() string }
type Dog struct{}

func (*Dog) Sound() string { return "woof" } // pointer receiver!

func main() {
	// A slice of the interface type can hold different concrete types together.
	// Each call such as s.Area() runs the method of whichever type is inside.
	shapes := []Shape{Circle{1}, Square{2}}
	for _, s := range shapes {
		fmt.Println(describe(s))
	}

	fmt.Println(countBytes(strings.NewReader("hello world"))) // any Reader works: 11

	// TYPE ASSERTION: `s.(Square)` asks "is the concrete value inside this interface a
	// Square?" The comma-ok form returns (value, true) on success or (zero value, false)
	// on failure. WITHOUT the `ok`, a wrong guess would panic.
	var s Shape = Square{3}
	if sq, ok := s.(Square); ok {
		fmt.Println("square side", sq.S) // we regained access to the field S
	}
	_, isCircle := s.(Circle)
	fmt.Println(isCircle) // false

	// You can also assert to ANOTHER INTERFACE to test for an optional capability.
	// Square has no String method, so this takes the else branch.
	if st, ok := s.(fmt.Stringer); ok {
		fmt.Println(st)
	} else {
		fmt.Println("Square is not a Stringer")
	}

	// THE EMPTY INTERFACE: `any` (same as interface{}) lists no methods, so every type
	// satisfies it. It can hold anything, but the compiler no longer knows what is
	// inside, so you must use a type assertion to get it back. Prefer real interfaces
	// or generics (lesson 14) when you can.
	var anything any = 42
	fmt.Println(anything.(int) + 1) // 43

	// METHOD SETS: which methods a type "has" depends on whether it is a value or pointer.
	//   *Dog (pointer) has all methods, including the pointer-receiver Sound.
	//   Dog  (value)   only has value-receiver methods.
	// So only &Dog{} satisfies Animal. Reason: a value stored inside an interface is
	// not addressable, so Go could not safely hand the method a pointer to it.
	var a Animal = &Dog{} // ok
	// var b Animal = Dog{} // compile error: Sound has pointer receiver
	fmt.Println(a.Sound())

	// GOTCHA: an interface value is a PAIR (type, value). It equals nil only when BOTH
	// parts are nil. Here `d` is a nil *Dog, but assigning it to `an` stores
	// (type=*Dog, value=nil): the type part is filled in, so `an != nil`.
	// Calling methods still works (the method is found via the type), but
	// `if err != nil` checks go wrong if you ever return a typed nil pointer as an
	// `error`. Fix: return a literal `nil`, not a nil pointer variable.
	var d *Dog // nil pointer
	var an Animal = d
	fmt.Println(d == nil, an == nil) // true false !!
}

// UNDER THE HOOD
// An interface value is TWO machine words (8 bytes each on 64-bit):
//
//	+-----------------+-----------------+
//	| itab / type ptr | data pointer    |
//	+-----------------+-----------------+
//
//   - The first word points to an "itab": a runtime table describing the concrete
//     type and holding the addresses of the methods that satisfy the interface.
//   - The second word points to the actual value (copied to the heap if needed).
//
// Calling s.Area() means: read the itab, find the Area function address, jump to it
// (dynamic dispatch). That costs a small extra step compared with a direct call,
// similar to a C++ virtual function call. Interface calls also usually prevent
// inlining, which is rarely a problem except in very hot loops.
