// Lesson 3: Control flow. Go has only `for` as a loop keyword.
//
// WHAT YOU WILL LEARN
//   - `if` with an init statement and how variable SCOPE works
//   - The one `for` keyword and its four shapes (C-style, while, infinite, range)
//   - `switch` (no fallthrough by default), tagless switch, and type switch
//   - Labels for breaking out of nested loops
//   - The Go 1.22 change that gives every loop iteration its own variable
//
// BACKGROUND: "control flow" means choosing which line of code runs next:
// branching (if/switch) and repeating (loops). By default code runs top to bottom.
package main

import "fmt"

func main() {
	// IF WITH AN INIT STATEMENT
	// Form: if <statement>; <condition> { ... }
	// `n := 7` runs first, then `n%2 == 0` is tested. (% is the remainder operator:
	// n%2 is 0 for even numbers.) The variable n exists only inside this if / else-if /
	// else chain: that region is its SCOPE. Outside the chain, `n` is undefined here.
	// Scoping variables tightly stops them leaking into code that should not touch them.
	// Note: no parentheses around the condition, but the braces are mandatory.
	if n := 7; n%2 == 0 {
		fmt.Println("even")
	} else if n > 5 {
		fmt.Println("odd and > 5")
	}

	// THE `for` LOOP: Go has no `while` or `do`. One keyword covers everything.
	// Shape 1: C-style `for init; condition; post`.
	//   init runs once; condition is checked before every pass; post runs after every pass.
	for i := 0; i < 3; i++ { // C-style
		fmt.Print(i, " ") // Print (no "ln") adds no newline; we add it ourselves below
	}
	fmt.Println()

	// Shape 2: only a condition. This is what other languages call `while`.
	n := 0
	for n < 3 { // while-style
		n++ // n++ means n = n + 1 (a statement in Go, not an expression)
	}
	// Shape 3: no condition at all means "forever"; leave with `break`
	// (or `return`, or a panic).
	for { // infinite
		n++
		if n > 5 {
			break // exit the innermost loop immediately
		}
	}
	fmt.Println(n) // 6

	// Shape 4: `range` walks over a collection. It can iterate slices, arrays, maps,
	// strings, channels, integers, and (Go 1.23+) iterator functions.
	// For a slice it yields (index, copy of the element) on each pass.
	for i, v := range []string{"a", "b"} {
		fmt.Println(i, v)
	}
	// Ranging over an integer N counts 0..N-1. (Go 1.22+)
	for i := range 3 { // range over int (Go 1.22+)
		fmt.Print(i, " ")
	}
	fmt.Println()
	// Ranging over a string decodes UTF-8 (lesson 7). Characters like é take MORE than
	// one byte, so `i` is the position in BYTES where each character starts (0, 1, 3),
	// not 0, 1, 2. `r` is the decoded character (a rune).
	for i, r := range "héy" { // strings iterate by rune; i is the BYTE offset
		fmt.Print(i, ":", string(r), " ")
	}
	fmt.Println()

	// SWITCH
	// Compares `day` against each case in order and runs the FIRST match.
	// Unlike C/JS there is no accidental "falling through" into the next case, so no
	// `break` is needed. A case can list several values separated by commas.
	switch day := 3; day {
	case 1, 7:
		fmt.Println("weekend")
	case 2, 3, 4, 5, 6:
		fmt.Println("weekday")
		fallthrough // opt in to falling through: also run the NEXT case body unconditionally
	default: // runs when no case matched (or via fallthrough, as here)
		fmt.Println("reached via fallthrough")
	}

	// A switch with no value after the keyword tests each case as a boolean condition,
	// top to bottom. It reads more cleanly than a long if / else-if chain.
	score := 85
	switch {
	case score >= 90:
		fmt.Println("A")
	case score >= 80:
		fmt.Println("B") // 85 matches here; the later cases are skipped
	}

	// TYPE SWITCH
	// `any` (alias of interface{}) can hold a value of ANY type (lesson 10). A value
	// inside it remembers its real ("dynamic") type at runtime. `v.(type)` is special
	// syntax valid only in a switch: it branches on that stored type, and inside each
	// case `t` automatically has that concrete type (so t+1 works for int, len(t) for string).
	for _, v := range []any{1, "hi", 2.5, nil, []int{1}} {
		switch t := v.(type) {
		case int:
			fmt.Println("int", t+1)
		case string:
			fmt.Println("string", len(t))
		case nil:
			fmt.Println("nil")
		default: // 2.5 and []int{1} land here; t keeps the original interface type
			fmt.Printf("other %T\n", t)
		}
	}

	// LABELS
	// A plain `break`/`continue` affects only the INNERMOST loop. Put a label (a name
	// followed by a colon) before a loop to target the OUTER one.
	//   continue outer -> skip to the next iteration of the outer loop
	//   break outer    -> leave both loops
	// Trace: i=0 prints j=0,1 then j==2 triggers `continue outer`;
	//        i=1 does the same; i=2 hits `break outer` on its first inner pass.
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

	// PER-ITERATION LOOP VARIABLES (Go 1.22+)
	// A CLOSURE is a function that remembers the variables around it (lesson 4).
	// Before Go 1.22 there was ONE `i` shared by the whole loop, so every closure saw
	// its final value (3 3 3), the same trap as JavaScript's `var` in loops.
	// Now each iteration gets a fresh copy of `i`, so each closure keeps its own value.
	// (This is selected by the `go 1.23` line in go.mod; older go.mod versions keep
	// the old behaviour for compatibility.)
	var fns []func() int
	for i := 0; i < 3; i++ {
		fns = append(fns, func() int { return i })
	}
	for _, f := range fns {
		fmt.Print(f(), " ") // 0 1 2
	}
	fmt.Println()
}
