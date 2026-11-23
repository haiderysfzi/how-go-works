package main

import "testing"

// A TEST is an ordinary function that checks that code behaves as expected.
// Rules for `go test` to find it:
//   - the file name ends in _test.go
//   - the function name starts with `Test` followed by a capital letter
//   - it takes one parameter, t *testing.T, your handle for reporting results
//
// t.Errorf records a failure and the test continues; t.Fatalf records a failure and
// stops the test at once. A test PASSES if it finishes without recording any failure.
// "got"/"want" naming is the Go convention for the actual and expected values.
func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Errorf("Add(2,3) = %d, want 5", got)
	}
}

// TABLE-DRIVEN TESTS are the Go idiom: describe the cases as DATA (a slice of structs)
// and run them all through one loop. Adding a case is one new line, not a new function.
// Each case runs as a named SUBTEST with t.Run, so a failure shows exactly which case
// broke (e.g. TestAbs/negative) and you can run one alone:
//
//	go test -run 'TestAbs/negative' ./20-testing
func TestAbs(t *testing.T) {
	// A slice of an anonymous struct type (lesson 8) declared and filled in place.
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"positive", 5, 5},
		{"negative", -5, 5},
		{"zero", 0, 0},
	}
	for _, tc := range tests { // tc = "test case"; Go 1.22 gives each iteration its own tc
		t.Run(tc.name, func(t *testing.T) {
			if got := Abs(tc.in); got != tc.want {
				t.Errorf("Abs(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// BENCHMARKS measure speed. Name starts with `Benchmark`, parameter is b *testing.B.
// Run with: go test -bench=. -benchmem   (-benchmem also reports memory allocations)
// The framework chooses b.N for you: it runs the loop with a small N, then keeps
// increasing N until the timing is stable, and prints the average time per iteration
// (ns/op). `for range b.N` is the Go 1.22+ spelling of "repeat b.N times".
func BenchmarkAdd(b *testing.B) {
	for range b.N {
		Add(1, 2)
	}
}

// FUZZING feeds a function many automatically generated, mutated inputs, looking for
// ones that break your assumptions or crash it. f.Add supplies "seed" inputs to start
// from. The function given to f.Fuzz is run for every generated value of its
// parameters (here one int `n`).
// With plain `go test` the fuzz function only runs the seeds as ordinary tests;
// real fuzzing needs the -fuzz flag: go test -fuzz=FuzzAbs -fuzztime=5s
// When it finds a failing input it saves it under testdata/ so it becomes a permanent
// regression test.
func FuzzAbs(f *testing.F) {
	f.Add(-3)
	f.Fuzz(func(t *testing.T, n int) {
		// The smallest int has no positive counterpart, so -n == n (it overflows). A
		// fuzzer will find that edge case on its own. Here we skip it. t.Skip marks
		// the input as "not interesting" instead of failing.
		if n != -n && Abs(n) < 0 { // MinInt overflows: fuzzing finds such edge cases
			t.Skip()
		}
	})
}

// RUNNING TESTS: handy flags
//   go test ./20-testing            run tests, print only failures
//   go test -v ./20-testing         verbose: list every test and subtest
//   go test -race ./20-testing      also detect data races (lesson 17)
//   go test -cover ./20-testing     report what percent of statements the tests executed
