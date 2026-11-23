// Lesson 14: Generics (type parameters), Go 1.18+.
//
// WHAT YOU WILL LEARN
//   - The problem generics solve
//   - Type parameters, constraints, and type inference
//   - Generic functions (Sum, Map, Filter, MaxOf) and generic types (Stack, Pair)
//   - How the compiler implements generics under the hood
//
// BACKGROUND: suppose you need a function that sums a slice. Without generics you must
// write SumInts([]int), SumFloats([]float64), ... with identical bodies, or give up
// type safety with `any`. GENERICS let you write the function once with a placeholder
// for the type (a TYPE PARAMETER, in square brackets) and have the compiler check it
// for every concrete type you actually use.
package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// A CONSTRAINT says which types a type parameter may be. It is written as an interface
// whose "type set" lists allowed types joined by `|`. The tilde in ~int means "int AND
// every type whose underlying type is int", such as `type MyInt int` below. Without the
// tilde, MyInt would be rejected. Constraint tells the compiler which operations (here
// `+`) are legal on T.
type Number interface {
	~int | ~int64 | ~float64
}

// [T Number] declares a type parameter named T constrained by Number. When you call
// Sum([]int{1,2,3}), the compiler INFERS T = int from the argument, so you need not
// write Sum[int](...). `var total T` starts at T's zero value.
func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x // allowed because every type in Number supports +
	}
	return total
}

// Generic functions may have several type parameters. `any` as the constraint means
// "no restriction at all". Map transforms each T into a U with the function f, e.g.
// ints to strings. (Not the key-value Map from lesson 6; this is the "map over a
// list" idea from functional programming.)
func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs)) // pre-sized: we know exactly how many results
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

// Filter keeps only the elements for which `keep` returns true.
func Filter[T any](xs []T, keep func(T) bool) []T {
	var out []T
	for _, x := range xs {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

// cmp.Ordered is a ready-made constraint from the standard library for every type that
// supports < and > (all integers, floats, and strings). That is why `a > b` compiles.
func MaxOf[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// GENERIC TYPES: a struct can also take type parameters. Stack[int] and Stack[string]
// are two distinct types made from one definition. The type parameter is repeated in
// the receiver: (s *Stack[T]).
type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// Pop returns the top item and a bool saying whether there was one. When the stack is
// empty we must return SOMETHING of type T, so we return `zero`: a variable declared
// with `var zero T`, which holds T's zero value whatever T turns out to be.
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1] // shrink the slice by one
	return v, true
}

// `comparable` is a built-in constraint for types usable with == (and thus as map keys,
// lesson 6). K must be comparable; V may be anything.
type Pair[K comparable, V any] struct {
	Key K
	Val V
}

type MyInt int

func main() {
	// Three different element types, one function. The last uses MyInt, which only
	// works because of the ~int in Number.
	fmt.Println(Sum([]int{1, 2, 3}), Sum([]float64{1.5, 2.5}), Sum([]MyInt{4, 5}))         // 6 4 9
	fmt.Println(Map([]int{1, 2, 3}, func(i int) string { return strings.Repeat("*", i) })) // [* ** ***]
	fmt.Println(Filter([]int{1, 2, 3, 4}, func(i int) bool { return i%2 == 0 }))           // [2 4]
	fmt.Println(MaxOf("apple", "pear"), MaxOf(3, 9))                                       // type inferred: pear 9

	// For generic TYPES you normally write the type argument explicitly: Stack[string].
	var s Stack[string]
	s.Push("a")
	s.Push("b")
	v, _ := s.Pop()
	fmt.Println(v) // b: last in, first out

	fmt.Println(Pair[string, int]{"age", 30}) // {age 30}

	// The standard library now offers generic helpers in the `slices` and `maps`
	// packages: no more hand-written loops for sorting or searching.
	// BinarySearch requires a sorted slice; it returns (position, found).
	xs := []int{3, 1, 2}
	slices.Sort(xs) // sorts in place: [1 2 3]
	idx, found := slices.BinarySearch(xs, 2)
	fmt.Println(xs, idx, found, slices.Contains(xs, 9), slices.Max(xs)) // [1 2 3] 1 true false 3
}

// UNDER THE HOOD: "GC shape stenciling".
// C++ templates generate a full copy of the function for each type (fast, but big
// binaries). Go instead compiles one copy per "GC shape": types with the same memory
// layout share code. All pointer types share one version; int and int64 may get
// their own. A hidden DICTIONARY (a table of type-specific info such as method
// addresses) is passed in at runtime so the shared code can behave correctly for each
// type. You get smaller binaries and faster builds, with a small runtime cost
// in some cases.
//
// ADVICE: do not reach for generics until you would otherwise copy-paste the same
// code for several types. For behaviour that differs by type, an interface (lesson 10)
// is usually clearer.
