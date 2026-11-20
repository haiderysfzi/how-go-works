// Lesson 14: Generics (type parameters), Go 1.18+.
package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// A constraint is an interface listing allowed types. ~int includes named types based on int.
type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

// Generic functions with multiple type params.
func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

func Filter[T any](xs []T, keep func(T) bool) []T {
	var out []T
	for _, x := range xs {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

func MaxOf[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Generic types.
type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

type Pair[K comparable, V any] struct {
	Key K
	Val V
}

type MyInt int

func main() {
	fmt.Println(Sum([]int{1, 2, 3}), Sum([]float64{1.5, 2.5}), Sum([]MyInt{4, 5}))
	fmt.Println(Map([]int{1, 2, 3}, func(i int) string { return strings.Repeat("*", i) }))
	fmt.Println(Filter([]int{1, 2, 3, 4}, func(i int) bool { return i%2 == 0 }))
	fmt.Println(MaxOf("apple", "pear"), MaxOf(3, 9)) // type inferred

	var s Stack[string]
	s.Push("a")
	s.Push("b")
	v, _ := s.Pop()
	fmt.Println(v)

	fmt.Println(Pair[string, int]{"age", 30})

	// The standard library now has generic helpers.
	xs := []int{3, 1, 2}
	slices.Sort(xs)
	idx, found := slices.BinarySearch(xs, 2)
	fmt.Println(xs, idx, found, slices.Contains(xs, 9), slices.Max(xs))
}

// Under the hood: Go uses "GC shape stenciling": one compiled copy per memory shape
// (for example all pointer types share one), plus a dictionary passed at runtime.
// Don't reach for generics until you'd otherwise copy-paste the same code.
