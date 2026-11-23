// Lesson 2: Variables, types, zero values and constants.
//
// WHAT YOU WILL LEARN
//   - The three ways to declare a variable, and when each is used
//   - Zero values: why a Go variable is never "undefined"
//   - The basic types and how many bits/bytes they use
//   - Why Go refuses to convert between types automatically
//   - Constants and `iota` (Go's way of doing enums)
//
// BACKGROUND: what is a "type"?
//
//	In memory everything is just bits. A TYPE tells the compiler how many bits a value
//	occupies and how to interpret them. The same 8 bits `01000001` are the number 65 if
//	read as an integer, or the letter 'A' if read as a character. Go checks types at
//	COMPILE time (it is "statically typed"), so you find mistakes before running.
package main

import "fmt"

// Package-level (global to this package) declarations live outside any function.
// A CONSTANT is a value fixed at compile time; it can never change while the program runs.
//
// `Pi` has no type written next to it, so it is an UNTYPED constant. An untyped
// constant is like a plain mathematical number that takes the type it is needed as:
// in `Pi * someFloat32` it becomes a float32, in `Pi * someFloat64` a float64.
const Pi = 3.14159 // untyped constant: adapts to float32/float64 as needed

// `type X Y` creates a NEW type X whose memory layout is the same as Y. Weekday is
// stored exactly like an int, but the compiler treats the two as different types, so
// you cannot accidentally pass a plain int where a Weekday is expected without saying so.
// This is how you give meaning to numbers ("this int is a day, not a count").
type Weekday int // a new named type based on int

// `iota` is a counter that exists only inside a const ( ... ) block. It starts at 0 on
// the first line and goes up by one on every following line. When a line omits the
// expression, it repeats the previous one, so this block means:
//
//	Sunday  Weekday = 0
//	Monday  Weekday = 1
//	Tuesday Weekday = 2
//
// Go has no `enum` keyword; named constants + iota fill that role.
const (
	Sunday Weekday = iota // iota counts 0, 1, 2... inside a const block
	Monday
	Tuesday
)

func main() {
	// THREE WAYS TO DECLARE A VARIABLE
	// A variable is a named box in memory that holds a value of one fixed type.
	var a int = 10 // explicit: keyword `var`, name, type, initial value
	var b = 20     // type inferred: the compiler sees the literal 20 and picks int
	c := 30        // short declaration `:=` means "declare AND assign, infer the type".
	//                It only works INSIDE functions. It is by far the most common form.
	fmt.Println(a, b, c)
	// Rule of thumb: use `:=` inside functions; use `var` when you want the zero value
	// (below) or need a package-level variable.

	// ZERO VALUES
	// In C, reading an uninitialised variable gives garbage; in JS you get `undefined`.
	// Go does neither: when you declare a variable without a value, Go sets all its
	// memory to zero bits, which means a sensible "empty" value for each type.
	var i int            // 0
	var f float64        // 0
	var s string         // "" (empty string)
	var ok bool          // false
	var p *int           // nil (a pointer pointing at nothing; see lesson 9)
	var sl []int         // nil slice (lesson 5)
	var m map[string]int // nil map (lesson 6)
	// %d prints an integer, %v prints "whatever is natural for this type",
	// %q prints a string wrapped in quotes (so you can SEE that it is empty).
	// `sl == nil` and `m == nil` are comparisons that produce a bool.
	fmt.Printf("%d %v %q %v %v %v %v\n", i, f, s, ok, p, sl == nil, m == nil)

	// BASIC TYPES
	// The number in a type name is how many BITS it uses. 8 bits = 1 byte.
	//   uint8  : unsigned (no negatives), 8 bits  -> 0..255
	//   int64  : signed, 64 bits                  -> about +-9.2 * 10^18
	//   float32: decimal number, 32 bits (about 7 significant digits)
	// Plain `int` and `uint` are 64 bits on 64-bit machines (32 on 32-bit machines).
	var (
		u8  uint8   = 255
		i64 int64   = 1 << 40 // `<<` shifts bits left: 1 << 40 means 2 to the power 40
		f32 float32 = 1.5
		r   rune    = 'G' // rune = int32, an alias meaning "one Unicode code point";
		//                    'G' in single quotes is a rune literal; it prints as 71
		by   byte = 'A'           // byte = uint8, an alias meaning "one raw byte"; 'A' is 65
		cplx      = complex(1, 2) // complex number 1+2i (rarely used)
	)
	fmt.Println(u8, i64, f32, r, by, cplx)

	// NO IMPLICIT CONVERSION
	// Many languages quietly turn an int into a float when you mix them. That hides
	// bugs (precision loss, overflow), so Go forces you to say what you mean.
	x := 7   // int
	y := 2.5 // float64 (the default type for a literal with a decimal point)
	// fmt.Println(x * y)           // compile error: mismatched types int and float64
	fmt.Println(float64(x) * y) // explicit conversion: T(value) makes a value of type T

	// `time` here is our own type defined at the bottom of this file (based on int).
	// Even though it is "just an int" inside, the untyped constant 3 is allowed to
	// become a `time`, but a variable of type int would NOT be assignable to it
	// without writing time(thatVariable).
	var d time = 3 // named types are distinct from their base type
	fmt.Println(d)

	// INTEGER OVERFLOW
	// A uint8 can only hold 0..255. 250 + 10 = 260, which needs 9 bits. The extra
	// bit is dropped, leaving 260 - 256 = 4. This "wraps around" like a car odometer,
	// with NO error or panic. Choose a type wide enough for your data.
	var small uint8 = 250
	small += 10
	fmt.Println(small) // 4

	// CONSTANTS
	fmt.Println(Pi*2, Monday, Tuesday) // Monday and Tuesday print as their numbers 1 and 2
	// Constants are evaluated by the compiler with ARBITRARY precision, so this 1 << 100
	// is legal even though no integer type holds 2^100. It must be reduced to something
	// that fits before it is stored in a variable. Here, 2^100 >> 98 = 2^2 = 4.
	const big = 1 << 100 // fine: constants have arbitrary precision
	fmt.Println(big >> 98)

	// TYPE INSPECTION: %T prints the type of a value, handy when you are unsure what
	// the compiler inferred. Expect: int float64 string main.Weekday
	fmt.Printf("%T %T %T %T\n", a, y, s, Monday)

	// MULTIPLE ASSIGNMENT: the right-hand side is fully evaluated first, THEN assigned,
	// so you can swap two variables with no temporary.
	m1, m2 := 1, 2
	m1, m2 = m2, m1
	fmt.Println(m1, m2) // 2 1

	// BLANK IDENTIFIER `_`: a write-only black hole. Go refuses unused variables, so
	// when a function returns a value you do not need, assign it to `_`.
	// fmt.Println returns (bytesWritten int, err error); we ignore the first one.
	_, err := fmt.Println("ignoring the byte count")
	_ = err // we also "use" err so the compiler does not complain about an unused variable
}

// A named type used above. It lives down here only to show that declaration order of
// package-level items does not matter in Go (unlike inside functions).
type time int
