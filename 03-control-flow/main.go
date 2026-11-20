// Lesson 3: Control flow. Go has only `for` as a loop keyword.
package main

import "fmt"

func main() {
	// if can have an init statement; the variable is scoped to the if/else chain.
	if n := 7; n%2 == 0 {
		fmt.Println("even")
	} else if n > 5 {
		fmt.Println("odd and > 5")
	}

	// The three flavours of for.
	for i := 0; i < 3; i++ { // C-style
		fmt.Print(i, " ")
	}
	fmt.Println()

	n := 0
	for n < 3 { // while-style
		n++
	}
	for { // infinite
		n++
		if n > 5 {
			break
		}
	}
	fmt.Println(n)

	// range iterates slices, arrays, maps, strings, channels, integers, functions.
	for i, v := range []string{"a", "b"} {
		fmt.Println(i, v)
	}
	for i := range 3 { // range over int (Go 1.22+)
		fmt.Print(i, " ")
	}
	fmt.Println()
	for i, r := range "héy" { // strings iterate by rune; i is the BYTE offset
		fmt.Print(i, ":", string(r), " ")
	}
	fmt.Println()

	// switch: no fallthrough by default, no break needed.
	switch day := 3; day {
	case 1, 7:
		fmt.Println("weekend")
	case 2, 3, 4, 5, 6:
		fmt.Println("weekday")
		fallthrough // opt in to falling through
	default:
		fmt.Println("reached via fallthrough")
	}

	// switch with no tag is a cleaner if/else chain.
	score := 85
	switch {
	case score >= 90:
		fmt.Println("A")
	case score >= 80:
		fmt.Println("B")
	}

	// Type switch: branch on the dynamic type of an interface.
	for _, v := range []any{1, "hi", 2.5, nil, []int{1}} {
		switch t := v.(type) {
		case int:
			fmt.Println("int", t+1)
		case string:
			fmt.Println("string", len(t))
		case nil:
			fmt.Println("nil")
		default:
			fmt.Printf("other %T\n", t)
		}
	}

	// Labels with break / continue.
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 2 {
				continue outer
			}
			if i == 2 {
				break outer
			}
			fmt.Print(i, j, " ")
		}
	}
	fmt.Println()

	// Go 1.22: each loop iteration gets its OWN copy of the variable,
	// so closures capture what you expect (the classic JS `var` gotcha is gone).
	var fns []func() int
	for i := 0; i < 3; i++ {
		fns = append(fns, func() int { return i })
	}
	for _, f := range fns {
		fmt.Print(f(), " ") // 0 1 2
	}
	fmt.Println()
}
