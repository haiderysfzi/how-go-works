// Lesson 7: Strings are immutable UTF-8 byte sequences. Runes are code points.
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	s := "héllo, 世界"

	fmt.Println(len(s))                    // bytes, not characters
	fmt.Println(utf8.RuneCountInString(s)) // characters (code points)
	fmt.Println(s[1], s[0:1])              // indexing gives a BYTE

	// Iterating with range decodes UTF-8 into runes.
	for i, r := range "aé世" {
		fmt.Printf("%d:%c(%d bytes) ", i, r, utf8.RuneLen(r))
	}
	fmt.Println()

	// Convert to []rune to index by character, []byte to mutate.
	runes := []rune(s)
	fmt.Println(string(runes[1]), len(runes))
	b := []byte("abc")
	b[0] = 'X'
	fmt.Println(string(b))
	// s[0] = 'X' // compile error: strings are immutable

	// strings package.
	fmt.Println(strings.ToUpper("go"), strings.Contains("seafood", "foo"),
		strings.Split("a,b,c", ","), strings.TrimSpace("  hi  "),
		strings.Fields(" a  b "), strings.Repeat("ab", 3),
		strings.ReplaceAll("a-b-c", "-", "+"), strings.HasPrefix("golang", "go"))

	// Building strings efficiently: += copies every time (quadratic).
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&sb, "%d,", i)
	}
	fmt.Println(sb.String())

	// Raw strings keep backslashes and newlines.
	fmt.Println(`C:\path\n stays literal`)

	// Formatting verbs.
	type P struct {
		Name string
		Age  int
	}
	p := P{"Ada", 36}
	fmt.Printf("%v | %+v | %#v | %T\n", p, p, p, p)
	fmt.Printf("%5d|%-5d|%05d|%x|%b|%8.3f|%q|%U\n", 42, 42, 42, 255, 5, 3.14159, "hi", '世')

	// Comparison is byte-wise and strings are comparable with ==.
	fmt.Println("a" < "b", "abc" == "ab"+"c")

	// Invalid UTF-8 becomes U+FFFD when ranged.
	for _, r := range "\xff" {
		fmt.Println(r == utf8.RuneError)
	}
}
