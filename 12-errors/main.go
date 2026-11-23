// Lesson 12: Errors are ordinary values, returned and checked explicitly.
//
// WHAT YOU WILL LEARN
//   - What the `error` type really is
//   - Sentinel errors and custom error types
//   - Wrapping errors with %w, and inspecting them with errors.Is / errors.As
//   - Joining several errors
//
// BACKGROUND: many languages use EXCEPTIONS: when something fails, execution jumps out
// of the current function to a distant catch block, invisibly. Go deliberately does not.
// A function that can fail returns an `error` as an extra value, and the caller must
// look at it. That makes every failure point visible in the code, at the cost of some
// repetition (`if err != nil { ... }`).
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// The built-in `error` type is just an interface with one method:
//
//	type error interface { Error() string }
//
// Any type with that method IS an error. nil means "no error".

// SENTINEL ERROR: a package-level variable created once with errors.New. Callers
// compare against it to learn WHAT went wrong. By convention these are named ErrXxx.
var ErrNotFound = errors.New("not found")

// CUSTOM ERROR TYPE: use a struct when the caller needs structured data (which field
// failed?) and not only a message. It holds the field name plus the underlying
// reason (Err). Pointer receivers: we return *ValidationError.
type ValidationError struct {
	Field string
	Err   error
}

func (v *ValidationError) Error() string { return "validation " + v.Field + ": " + v.Err.Error() }

// Unwrap returns the error that this one wraps. errors.Is and errors.As call it to walk
// down the chain of causes.
func (v *ValidationError) Unwrap() error { return v.Err } // enables errors.Is/As through the chain

func find(id int) error {
	if id < 0 {
		return &ValidationError{"id", errors.New("must be positive")}
	}
	// WRAPPING with fmt.Errorf and the %w verb: this builds a NEW error whose message
	// adds context ("find 5: not found") while keeping the original error inside,
	// so callers can still detect ErrNotFound. (%v instead of %w would keep only the text.)
	return fmt.Errorf("find %d: %w", id, ErrNotFound) // %w WRAPS, keeping the original
}

func main() {
	// THE IDIOMATIC PATTERN: call, check err immediately, handle or return it.
	// strconv.Atoi ("ASCII to integer") returns (int, error); "12x" is not a number.
	n, err := strconv.Atoi("12x")
	if err != nil {
		fmt.Println("error:", err)
	}
	_ = n

	// errors.Is(err, target): "is `target` anywhere in the chain of err?" It unwraps
	// step by step and compares each link with ==. A plain `err == ErrNotFound` would
	// be FALSE here because err is the wrapper, not the sentinel itself.
	err = find(5)
	fmt.Println(err, "|", errors.Is(err, ErrNotFound)) // find 5: not found | true

	// errors.As(err, &target): "is there an error of THIS TYPE in the chain?" If so, it
	// stores it in `target` so you can read its fields. We pass the ADDRESS of `ve`
	// (a pointer to a pointer variable) so As can fill it in.
	err = find(-1)
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Println("field:", ve.Field) // id
	}

	// errors.Join (Go 1.20+) bundles several independent errors into one value.
	// Its message lists each error on its own line.
	joined := errors.Join(errors.New("a failed"), errors.New("b failed"))
	fmt.Println(joined)

	// Standard library errors work the same way: os.Open wraps the OS error, and
	// os.ErrNotExist is a sentinel that errors.Is can detect on any operating system.
	_, err = os.Open("/definitely/missing")
	fmt.Println(errors.Is(err, os.ErrNotExist)) // true

	// GUIDELINES
	// - Do not ignore errors with `_` unless you can justify it.
	// - Handle an error ONCE: either log it or return it, not both (otherwise the same
	//   failure appears in the logs many times).
	// - Add context as you return up the stack: fmt.Errorf("loading config: %w", err).
	// There are no exceptions: panics (next lesson) are only for programmer bugs.
}
