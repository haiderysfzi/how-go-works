// Lesson 6: Maps (hash tables).
package main

import (
	"fmt"
	"maps"
	"slices"
	"sort"
)

func main() {
	ages := map[string]int{"ada": 36, "bob": 25}
	ages["cy"] = 41

	// Reading a missing key returns the zero value, so use the "comma ok" form.
	fmt.Println(ages["nobody"]) // 0
	if v, ok := ages["bob"]; ok {
		fmt.Println("found", v)
	}

	delete(ages, "bob")
	fmt.Println(len(ages))

	// Iteration order is RANDOMISED on purpose. Never rely on it.
	// Sort the keys when you need a stable order.
	keys := make([]string, 0, len(ages))
	for k := range ages {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println(keys)
	fmt.Println(slices.Sorted(maps.Keys(ages))) // Go 1.23+ shortcut

	// Writing to a nil map panics; reading is fine. Use make.
	var nilMap map[string]int
	fmt.Println(nilMap["x"])
	// nilMap["x"] = 1 // panic: assignment to entry in nil map

	// Keys must be comparable (no slices, maps or functions). Structs and arrays work.
	type point struct{ x, y int }
	visited := map[point]bool{{1, 2}: true}
	fmt.Println(visited[point{1, 2}], visited[point{0, 0}])

	// Set idiom: map[T]struct{} (the empty struct uses zero bytes).
	set := map[string]struct{}{"a": {}, "b": {}}
	_, has := set["a"]
	fmt.Println(has)

	// Maps are reference-like: passing one lets the callee mutate it.
	mutate(ages)
	fmt.Println(ages["added"])

	// Not safe for concurrent writes! Use sync.Mutex or sync.Map (see lesson 17).

	// Counting words.
	counts := map[rune]int{}
	for _, r := range "hello" {
		counts[r]++ // missing key starts at 0
	}
	fmt.Println(counts['l'])

	clear(counts) // Go 1.21+: remove all entries
	fmt.Println(len(counts))
}

func mutate(m map[string]int) { m["added"] = 1 }
