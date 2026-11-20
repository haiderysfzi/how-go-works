// Lesson 2: Variables, types, zero values and constants.
package main

import "fmt"

// Package-level variables and constants.
const Pi = 3.14159 // untyped constant: adapts to float32/float64 as needed

type Weekday int // a new named type based on int

const (
	Sunday Weekday = iota // iota counts 0, 1, 2... inside a const block
	Monday
	Tuesday
)

func main() {
	// Three ways to declare:
	var a int = 10 // explicit
	var b = 20     // type inferred (int)
	c := 30        // short declaration, only inside functions
	fmt.Println(a, b, c)

	// Zero values: every variable is initialised, there is no "undefined".
	var i int
	var f float64
	var s string
	var ok bool
	var p *int
	var sl []int
	var m map[string]int
	fmt.Printf("%d %v %q %v %v %v %v\n", i, f, s, ok, p, sl == nil, m == nil)

	// Basic types.
	var (
		u8   uint8   = 255
		i64  int64   = 1 << 40
		f32  float32 = 1.5
		r    rune    = 'G' // rune = int32, a Unicode code point
		by   byte    = 'A' // byte = uint8
		cplx         = complex(1, 2)
	)
	fmt.Println(u8, i64, f32, r, by, cplx)

	// Go NEVER converts types implicitly.
	x := 7
	y := 2.5
	// fmt.Println(x * y)           // compile error: mismatched types int and float64
	fmt.Println(float64(x) * y) // explicit conversion
	var d time = 3              // named types are distinct from their base type
	fmt.Println(d)

	// Integer overflow wraps silently.
	var small uint8 = 250
	small += 10
	fmt.Println(small) // 4

	// Constants.
	fmt.Println(Pi*2, Monday, Tuesday)
	const big = 1 << 100 // fine: constants have arbitrary precision
	fmt.Println(big >> 98)

	// Type inspection.
	fmt.Printf("%T %T %T %T\n", a, y, s, Monday)

	// Multiple assignment and swap.
	m1, m2 := 1, 2
	m1, m2 = m2, m1
	fmt.Println(m1, m2)

	// Blank identifier discards a value.
	_, err := fmt.Println("ignoring the byte count")
	_ = err
}

type time int
