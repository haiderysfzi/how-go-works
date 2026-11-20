// Lesson 1: How a Go program is built and run.
// Run with: go run ./01-hello-and-compilation
package main // "main" is the special package that produces an executable

import (
	"fmt"
	"runtime"
)

// init runs automatically before main, once per package.
func init() {
	fmt.Println("init: runs before main")
}

// main is the entry point.
func main() {
	fmt.Println("Hello, Go!")
	fmt.Println("OS/Arch:", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Go version:", runtime.Version())
	fmt.Println("CPUs:", runtime.NumCPU())
}

// What happens when you run `go build`:
// 1. The toolchain parses and type-checks ALL your code (errors stop the build here,
//    including unused variables and unused imports).
// 2. It compiles each package to machine code and links everything, including the
//    Go runtime (scheduler, garbage collector), into ONE static binary.
// 3. You can copy that binary to another machine with the same OS/arch and run it
//    with no interpreter or VM installed.
//
// Compare with JavaScript: no JIT warm-up, no interpreter. Compiled ahead of time.
// Useful commands:
//   go run .            compile to a temp binary and run it
//   go build            produce a binary
//   go vet ./...        static checks for suspicious code
//   gofmt -w .          the one official code style (no style debates)
//   GOOS=linux GOARCH=amd64 go build   cross-compile with one env var
