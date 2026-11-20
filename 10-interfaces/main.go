// Lesson 10: Interfaces are satisfied IMPLICITLY (structural typing).
package main

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// An interface is a set of method signatures. No "implements" keyword.
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

// Accept interfaces, return concrete types.
func describe(s Shape) string {
	return fmt.Sprintf("%T area=%.2f perimeter=%.2f", s, s.Area(), s.Perimeter())
}

// Small interfaces are best. io.Reader has ONE method and powers the whole stdlib.
func countBytes(r io.Reader) int {
	n, _ := io.Copy(io.Discard, r)
	return int(n)
}

type Animal interface{ Sound() string }
type Dog struct{}

func (*Dog) Sound() string { return "woof" } // pointer receiver!

func main() {
	shapes := []Shape{Circle{1}, Square{2}}
	for _, s := range shapes {
		fmt.Println(describe(s))
	}

	fmt.Println(countBytes(strings.NewReader("hello world"))) // any Reader works

	// Type assertion: get the concrete value back out.
	var s Shape = Square{3}
	if sq, ok := s.(Square); ok {
		fmt.Println("square side", sq.S)
	}
	_, isCircle := s.(Circle)
	fmt.Println(isCircle)

	// Check for an optional extra capability.
	if st, ok := s.(fmt.Stringer); ok {
		fmt.Println(st)
	} else {
		fmt.Println("Square is not a Stringer")
	}

	// The empty interface (`any`) holds anything; you lose type safety.
	var anything any = 42
	fmt.Println(anything.(int) + 1)

	// Method sets: *Dog has Sound, plain Dog does not.
	var a Animal = &Dog{} // ok
	// var b Animal = Dog{} // compile error: Sound has pointer receiver
	fmt.Println(a.Sound())

	// GOTCHA: an interface is (type, value). It is nil only when BOTH are nil.
	var d *Dog // nil pointer
	var an Animal = d
	fmt.Println(d == nil, an == nil) // true false !!
}

// Under the hood: an interface value is two words, a pointer to the type's method table
// (itab) and a pointer to the data. Calling a method looks it up through the itab,
// so it is a small dynamic dispatch cost, comparable to a C++ virtual call.
