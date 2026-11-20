// Lesson 8: Structs and methods. Go has no classes.
package main

import "fmt"

type Rect struct {
	Width, Height float64 // exported (capitalised) fields are visible to other packages
	label         string  // lowercase = package-private
}

// Value receiver: works on a COPY.
func (r Rect) Area() float64 { return r.Width * r.Height }

// Pointer receiver: can modify the original and avoids copying big structs.
func (r *Rect) Scale(f float64) {
	r.Width *= f
	r.Height *= f
}

// Constructor by convention (there are no real constructors).
func NewRect(w, h float64) *Rect { return &Rect{Width: w, Height: h, label: "rect"} }

// String satisfies fmt.Stringer, so Println uses it.
func (r Rect) String() string { return fmt.Sprintf("Rect(%gx%g)", r.Width, r.Height) }

type Celsius float64 // methods can be attached to ANY named type you define

func (c Celsius) Fahrenheit() float64 { return float64(c)*9/5 + 32 }

func main() {
	r := Rect{Width: 3, Height: 4} // named fields (preferred)
	fmt.Println(r.Area())
	r.Scale(2) // Go takes &r automatically for addressable values
	fmt.Println(r, r.Area())

	p := NewRect(1, 1)
	p.Scale(5) // pointers auto-dereference too
	fmt.Println(*p)

	// Structs are values: assignment copies.
	a := Rect{Width: 1}
	b := a
	b.Width = 9
	fmt.Println(a.Width, b.Width)

	// Comparable if all fields are comparable.
	fmt.Println(Rect{Width: 1} == Rect{Width: 1})

	// Anonymous structs and struct tags.
	cfg := struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}{"localhost", 8080}
	fmt.Printf("%+v\n", cfg)

	fmt.Println(Celsius(100).Fahrenheit())

	// Method values and expressions.
	area := r.Area  // bound to r
	fn := Rect.Area // takes the receiver as first argument
	fmt.Println(area(), fn(Rect{Width: 2, Height: 2}))

	// Zero value is useful: design types so the zero value works (like sync.Mutex).
	var empty Rect
	fmt.Println(empty.Area())
}
