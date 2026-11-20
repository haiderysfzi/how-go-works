// Lesson 12: Errors are ordinary values, returned and checked explicitly.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// The error interface is just: Error() string

// Sentinel errors for comparison.
var ErrNotFound = errors.New("not found")

// Custom error type to carry data.
type ValidationError struct {
	Field string
	Err   error
}

func (v *ValidationError) Error() string { return "validation " + v.Field + ": " + v.Err.Error() }
func (v *ValidationError) Unwrap() error { return v.Err } // enables errors.Is/As through the chain

func find(id int) error {
	if id < 0 {
		return &ValidationError{"id", errors.New("must be positive")}
	}
	return fmt.Errorf("find %d: %w", id, ErrNotFound) // %w WRAPS, keeping the original
}

func main() {
	// The idiomatic pattern.
	n, err := strconv.Atoi("12x")
	if err != nil {
		fmt.Println("error:", err)
	}
	_ = n

	// errors.Is walks the wrap chain comparing values.
	err = find(5)
	fmt.Println(err, "|", errors.Is(err, ErrNotFound))

	// errors.As finds the first error of a given TYPE in the chain.
	err = find(-1)
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Println("field:", ve.Field)
	}

	// Combine independent errors.
	joined := errors.Join(errors.New("a failed"), errors.New("b failed"))
	fmt.Println(joined)

	// Standard library errors are inspectable too.
	_, err = os.Open("/definitely/missing")
	fmt.Println(errors.Is(err, os.ErrNotExist))

	// Don't ignore errors; don't both log AND return them (handle once).
	// There are no exceptions: panics (next lesson) are only for programmer bugs.
}
