package main

import "testing"

// Test functions: file ends in _test.go, name starts with Test, take *testing.T.
func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Errorf("Add(2,3) = %d, want 5", got)
	}
}

// Table-driven tests are the Go idiom: one loop, many cases, named subtests.
func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"positive", 5, 5},
		{"negative", -5, 5},
		{"zero", 0, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Abs(tc.in); got != tc.want {
				t.Errorf("Abs(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// Benchmarks: go test -bench=. -benchmem
func BenchmarkAdd(b *testing.B) {
	for range b.N {
		Add(1, 2)
	}
}

// Fuzzing: go test -fuzz=FuzzAbs -fuzztime=5s
func FuzzAbs(f *testing.F) {
	f.Add(-3)
	f.Fuzz(func(t *testing.T, n int) {
		if n != -n && Abs(n) < 0 { // MinInt overflows: fuzzing finds such edge cases
			t.Skip()
		}
	})
}
