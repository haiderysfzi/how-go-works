// Lesson 6: Maps (hash tables).
//
// WHAT YOU WILL LEARN
//   - What a map is and how lookups work internally
//   - Reading missing keys and the "comma ok" idiom
//   - Why iteration order is random
//   - Nil maps, key requirements, sets, counting, and concurrency limits
//
// BACKGROUND: a map stores KEY -> VALUE pairs (like a JS object/Map or a Python dict).
// Internally it is a HASH TABLE: Go runs the key through a HASH FUNCTION that turns it
// into a number, and uses that number to pick a "bucket" (a small group of slots) in an
// array. Looking a key up means: hash it, jump straight to its bucket, compare a few
// entries. That takes roughly constant time no matter how many entries exist.
// Because the position depends on the hash, the entries have no meaningful order.
package main

import (
	"fmt"
	"maps"
	"slices"
	"sort"
)

func main() {
	// A map type is written map[KeyType]ValueType. This literal creates and fills one.
	ages := map[string]int{"ada": 36, "bob": 25}
	ages["cy"] = 41 // assignment inserts a new key, or overwrites an existing one

	// READING A MISSING KEY returns the ZERO VALUE of the value type (0 for int); it is
	// NOT an error. That makes "is the key absent?" and "is the value 0?" look the same.
	// The "comma ok" form fixes that: the second result `ok` is true only if the key
	// really exists.
	fmt.Println(ages["nobody"]) // 0
	if v, ok := ages["bob"]; ok {
		fmt.Println("found", v)
	}

	// delete(map, key) removes a key; deleting a missing key is a harmless no-op.
	delete(ages, "bob")
	fmt.Println(len(ages)) // 2: ada and cy

	// ITERATION ORDER IS RANDOMISED ON PURPOSE.
	// The Go runtime deliberately starts each `range` at a random position so that
	// programmers never accidentally depend on an order that is not guaranteed.
	// If you need a stable order, collect the keys, sort them, and loop over the
	// sorted keys.
	keys := make([]string, 0, len(ages)) // pre-sized slice (lesson 5)
	for k := range ages {                // with one variable, range gives just the keys
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println(keys)
	// Go 1.23+ shortcut: maps.Keys returns an iterator over the keys and slices.Sorted
	// collects and sorts it in one go.
	fmt.Println(slices.Sorted(maps.Keys(ages)))

	// NIL MAP: `var m map[K]V` gives a nil map. Reading from it is fine (always "missing"),
	// but WRITING panics, because there is no hash table allocated yet.
	// Create a usable map with make(map[K]V) or a literal.
	var nilMap map[string]int
	fmt.Println(nilMap["x"])
	// nilMap["x"] = 1 // panic: assignment to entry in nil map

	// KEYS MUST BE COMPARABLE: the map needs `==` to tell two keys apart after hashing.
	// Slices, maps and functions cannot be compared, so they cannot be keys.
	// Numbers, strings, booleans, pointers, arrays and structs of those all can.
	// A struct key matches only if ALL its fields are equal.
	type point struct{ x, y int }
	visited := map[point]bool{{1, 2}: true}                 // {1, 2} is shorthand for point{1, 2}
	fmt.Println(visited[point{1, 2}], visited[point{0, 0}]) // true false

	// SET IDIOM: Go has no set type. Use a map and ignore the value. struct{} is the
	// empty struct, which takes ZERO bytes, so the map stores only the keys.
	// `struct{}{}` (abbreviated `{}` in the literal) is the one value of that type.
	set := map[string]struct{}{"a": {}, "b": {}}
	_, has := set["a"]
	fmt.Println(has) // true

	// MAPS ARE "REFERENCE-LIKE". The map variable holds a small pointer to the runtime's
	// hash table. Passing the map to a function copies only that pointer (Go passes
	// everything by value), so the function and the caller reach the SAME table.
	mutate(ages)
	fmt.Println(ages["added"]) // 1

	// NOT SAFE FOR CONCURRENT WRITES! If two goroutines write at the same time, the
	// runtime detects it and crashes the program ("fatal error: concurrent map writes").
	// Protect the map with a sync.Mutex or use sync.Map (lesson 17).

	// COUNTING: because a missing key reads as 0, `counts[r]++` works even the first
	// time. It means "read the current value (0 if absent), add 1, store it".
	counts := map[rune]int{}
	for _, r := range "hello" { // iterating a string yields runes (lesson 7)
		counts[r]++ // missing key starts at 0
	}
	fmt.Println(counts['l']) // 2

	clear(counts) // Go 1.21+: remove all entries but keep the map itself usable
	fmt.Println(len(counts))
}

func mutate(m map[string]int) { m["added"] = 1 }
