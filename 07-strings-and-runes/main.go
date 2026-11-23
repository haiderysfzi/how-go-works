// Lesson 7: Strings are immutable UTF-8 byte sequences. Runes are code points.
//
// WHAT YOU WILL LEARN
//   - How text is stored in memory (bytes, Unicode, UTF-8)
//   - The difference between a byte, a rune and a "character"
//   - Why len("世") is 3, and what `range` over a string does
//   - Efficient string building and the fmt formatting verbs
//
// BACKGROUND: computers only store numbers, so text needs an agreed mapping.
//   - UNICODE assigns every character in every language a unique number called a
//     CODE POINT: 'A' = 65 (written U+0041), 'é' = 233 (U+00E9), '世' = 19990 (U+4E16).
//   - UTF-8 is the way those code points are written as BYTES. Small code points use
//     1 byte (all of ASCII); larger ones use 2, 3 or 4 bytes:
//     'a'  -> 1 byte    'é' -> 2 bytes    '世' -> 3 bytes    '😀' -> 4 bytes
//   - A Go `string` is just a read-only sequence of BYTES (normally UTF-8). Internally
//     it is a 2-field header: {pointer to bytes, length in bytes}.
//   - A `rune` is Go's name for one Unicode code point (it is an alias of int32).
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	s := "héllo, 世界"
	// Byte layout: h(1) é(2) l(1) l(1) o(1) ,(1) space(1) 世(3) 界(3) = 14 bytes,
	// but only 9 characters.

	fmt.Println(len(s))                    // 14: len counts BYTES, not characters
	fmt.Println(utf8.RuneCountInString(s)) // 9: decodes UTF-8 and counts code points
	// Indexing returns ONE BYTE (a uint8 number). s[1] is the first byte of 'é', which is
	// not a whole character. s[0:1] is a slicing operation and produces a 1-byte string "h".
	fmt.Println(s[1], s[0:1])

	// ITERATING WITH range: Go decodes UTF-8 for you. On each pass you get:
	//   i = the BYTE OFFSET where the character begins
	//   r = the decoded rune (code point)
	// so for "aé世" the offsets are 0, 1, 3 (not 0, 1, 2).
	// %c prints a rune as a character; utf8.RuneLen tells how many bytes it takes.
	for i, r := range "aé世" {
		fmt.Printf("%d:%c(%d bytes) ", i, r, utf8.RuneLen(r))
	}
	fmt.Println()

	// CONVERSIONS
	// []rune(s) decodes the whole string into a slice of code points. Now runes[1] is
	// really the 2nd CHARACTER. (This allocates and copies, so do it only when needed.)
	// []byte(s) makes a mutable COPY of the bytes. Modify it, then convert back with string().
	runes := []rune(s)
	fmt.Println(string(runes[1]), len(runes)) // é 9
	b := []byte("abc")
	b[0] = 'X'
	fmt.Println(string(b)) // Xbc
	// s[0] = 'X' // compile error: strings are immutable. Because they never change, Go
	//            // can safely let many variables share the same bytes without copying.

	// The `strings` package holds the common text helpers. None modify the input;
	// each returns a new string (or slice of strings).
	fmt.Println(strings.ToUpper("go"), strings.Contains("seafood", "foo"),
		strings.Split("a,b,c", ","), strings.TrimSpace("  hi  "),
		strings.Fields(" a  b "), // Fields splits on runs of whitespace -> [a b]
		strings.Repeat("ab", 3), strings.ReplaceAll("a-b-c", "-", "+"),
		strings.HasPrefix("golang", "go"))

	// BUILDING STRINGS EFFICIENTLY
	// Because strings are immutable, `result += "x"` creates a brand-new string and copies
	// all the old bytes on EVERY step. Doing that in a loop is "quadratic": 10x more
	// pieces means ~100x more copying. strings.Builder keeps one growing byte buffer
	// and only produces the final string once.
	// Fprintf writes formatted text to anything with a Write method (an io.Writer,
	// lesson 10), and *Builder is one.
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&sb, "%d,", i)
	}
	fmt.Println(sb.String()) // 0,1,2,

	// RAW STRINGS use backticks: no escape sequences are processed, newlines are
	// allowed, and backslashes are literal. Ideal for regular expressions, Windows paths, SQL.
	fmt.Println(`C:\path\n stays literal`)

	// FORMATTING VERBS (the % codes in Printf)
	type P struct {
		Name string
		Age  int
	}
	p := P{"Ada", 36}
	// %v default format | %+v includes field names | %#v Go-syntax | %T the type
	fmt.Printf("%v | %+v | %#v | %T\n", p, p, p, p)
	// %5d  integer padded to width 5, right-aligned    %-5d left-aligned
	// %05d zero-padded                                  %x hexadecimal
	// %b   binary                                       %8.3f float, width 8, 3 decimals
	// %q   quoted string                                %U Unicode code point (U+4E16)
	fmt.Printf("%5d|%-5d|%05d|%x|%b|%8.3f|%q|%U\n", 42, 42, 42, 255, 5, 3.14159, "hi", '世')

	// COMPARISON: strings compare byte by byte, in dictionary-like order of the raw
	// bytes. `==` checks both length and content; `<` finds the first differing byte.
	fmt.Println("a" < "b", "abc" == "ab"+"c") // true true

	// INVALID UTF-8: the byte 0xff can never appear in valid UTF-8. Instead of failing,
	// `range` yields the "replacement character" U+FFFD (utf8.RuneError) and moves on
	// one byte. Go never refuses to hold bytes it cannot decode.
	for _, r := range "\xff" {
		fmt.Println(r == utf8.RuneError) // true
	}
}
