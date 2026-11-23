// Lesson 8: Structs and methods. Go has no classes.
//
// WHAT YOU WILL LEARN
//   - What a struct is and how it sits in memory
//   - Methods, and the difference between value and pointer receivers
//   - Exported vs unexported names (Go's version of public/private)
//   - Constructors by convention, Stringer, method values, struct tags
//
// BACKGROUND: a STRUCT groups several named values (FIELDS) into one value. In memory
// the fields are laid out one after another in a single block, so Rect below is simply
// two float64s (8 bytes each) followed by a string header: a contiguous chunk.
// Go has no classes or inheritance. You combine structs + methods + interfaces
// (lesson 10) + embedding (lesson 11) instead.
package main

import "fmt"

type Rect struct {
	// Visibility rule for ALL names in Go: starts with a CAPITAL letter = exported
	// (visible to other packages); lowercase = private to this package. There is no
	// `public`/`private` keyword.
	Width, Height float64 // exported (capitalised) fields are visible to other packages
	label         string  // lowercase = package-private
}

// METHODS
// A method is a function with an extra parameter before the name, called the RECEIVER.
// `(r Rect)` says "this function belongs to Rect, and inside it the value is called r".
// It is called as r.Area(). The receiver plays the role of `this`/`self` elsewhere, but
// it is an ordinary named parameter.
//
// VALUE receiver (r Rect): the method gets a COPY of the struct. Anything it changes
// is lost when it returns. Fine for read-only methods on small structs.
func (r Rect) Area() float64 { return r.Width * r.Height }

// POINTER receiver (r *Rect): the method gets the ADDRESS of the original struct (see
// lesson 9), so it can modify the caller's value, and large structs are not copied
// on every call.
// `r.Width *= f` is shorthand for `(*r).Width = (*r).Width * f`; Go inserts the
// dereference for you.
func (r *Rect) Scale(f float64) {
	r.Width *= f
	r.Height *= f
}

// Go has no constructors. By convention you write a plain function named NewTypeName
// that builds and returns a value. `&Rect{...}` creates the struct and returns a pointer to it.
// Here we set the private `label` field, which outside packages could not do directly.
func NewRect(w, h float64) *Rect { return &Rect{Width: w, Height: h, label: "rect"} }

// String has a special meaning: fmt.Stringer is the interface `interface{ String() string }`.
// fmt.Println checks whether a value has that method and, if so, uses it to produce
// the text. This is "implicit interface satisfaction" (lesson 10): we never wrote
// "implements Stringer".
func (r Rect) String() string { return fmt.Sprintf("Rect(%gx%g)", r.Width, r.Height) }

// Methods can be attached to ANY named type you define in your own package, not only
// structs. Celsius is a float64 underneath but now carries its own behaviour.
type Celsius float64 // methods can be attached to ANY named type you define

func (c Celsius) Fahrenheit() float64 { return float64(c)*9/5 + 32 }

func main() {
	// A STRUCT LITERAL. Naming the fields is preferred: it stays correct if fields are
	// reordered and unlisted fields get their zero value (label is "" here).
	r := Rect{Width: 3, Height: 4} // named fields (preferred)
	fmt.Println(r.Area())          // 12

	// Scale needs *Rect but `r` is a Rect. Because `r` is a variable (it has an address),
	// Go automatically rewrites r.Scale(2) into (&r).Scale(2).
	r.Scale(2)               // Go takes &r automatically for addressable values
	fmt.Println(r, r.Area()) // uses our String method: Rect(6x8) 48

	p := NewRect(1, 1) // p is a *Rect
	p.Scale(5)         // pointers auto-dereference too: p.Area() works on a *Rect as well
	fmt.Println(*p)    // *p is the Rect that p points to

	// STRUCTS ARE VALUES: assignment copies every field, so a and b are independent.
	a := Rect{Width: 1}
	b := a
	b.Width = 9
	fmt.Println(a.Width, b.Width) // 1 9

	// A struct is comparable with == if every field is comparable; it checks each field.
	fmt.Println(Rect{Width: 1} == Rect{Width: 1}) // true

	// ANONYMOUS STRUCT: a struct type declared on the spot with no name; handy for
	// one-off data such as test cases or config. The backtick text after each field
	// is a STRUCT TAG: metadata (just a string) that libraries read at runtime with
	// reflection. encoding/json uses `json:"host"` to pick the key name in JSON (lesson 21).
	cfg := struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}{"localhost", 8080}
	fmt.Printf("%+v\n", cfg) // %+v prints field names

	fmt.Println(Celsius(100).Fahrenheit()) // 212

	// METHOD VALUES AND EXPRESSIONS
	// r.Area (without parentheses) is a METHOD VALUE: a function with the receiver
	// already attached. Because Area has a value receiver, it captures a COPY of r now.
	// Rect.Area is a METHOD EXPRESSION: a plain function whose first argument is the
	// receiver, so you pass it explicitly.
	area := r.Area                                     // bound to r
	fn := Rect.Area                                    // takes the receiver as first argument
	fmt.Println(area(), fn(Rect{Width: 2, Height: 2})) // 48 4

	// USEFUL ZERO VALUE: `var empty Rect` is already valid (all zeros) with no setup.
	// Good Go types are designed so that this "just works" (sync.Mutex and bytes.Buffer do).
	var empty Rect
	fmt.Println(empty.Area()) // 0
}
