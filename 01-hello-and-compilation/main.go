// Lesson 1: How a Go program is built and run.
// Run with: go run ./01-hello-and-compilation
//
// WHAT YOU WILL LEARN
//   - What a "package" is and why `main` is special
//   - What `init` and `main` are, and the order they run in
//   - What really happens between typing `go build` and seeing your program run
//
// HOW TO READ THESE LESSONS
//   Read the comments top to bottom, run the program, compare the output with what
//   the comments predict, then CHANGE the code and run it again. Breaking things on
//   purpose (and reading the error message) is the fastest way to learn.

// Every Go file starts with a `package` line. A PACKAGE is a folder of .go files that
// are compiled together and share the same names. It is Go's unit of code organisation
// (similar to a module in JS/Python or a namespace in C#).
//
// The package called `main` is special: it tells the compiler "build this into a
// program you can run" instead of "build this into a library other code imports".
// A `main` package MUST contain a function named `main`; that is where execution begins.
package main // "main" is the special package that produces an executable

// `import` lists the other packages this file uses. Each name is a path into either
// the standard library (like "fmt") or a third-party module.
//
// Go is strict here: if you import a package and never use it, the build FAILS.
// This keeps dependency lists honest (no "I think we still use that?").
import (
	"fmt"     // "format": printing and string formatting (Println, Printf, Sprintf...)
	"runtime" // lets a program ask about the Go runtime itself: OS, CPU count, version...
)

// init is a special function name. Rules:
//   - You never call it yourself; Go calls it automatically.
//   - It runs ONCE per package, after all package-level variables are initialised
//     and BEFORE main starts.
//   - A package may even have several init functions; they run in the order written.
//
// Typical use: set up something that must exist before the program starts
// (register a database driver, validate configuration). Use sparingly, because
// hidden "runs before main" behaviour is hard to follow in large programs.
func init() {
	fmt.Println("init: runs before main")
}

// main is the entry point. When it returns, the whole program exits immediately,
// even if other goroutines (lesson 15) are still running.
//
// Order of events when you run this program:
//  1. The Go runtime starts up (sets up memory, the garbage collector, the scheduler).
//  2. Package-level variables are initialised, then init() runs.
//  3. main() runs.
//  4. main() returns -> the process exits with status 0.
func main() {
	// fmt.Println prints its arguments separated by spaces, then a newline.
	// The capital P matters: in Go, names starting with an UPPERCASE letter are
	// "exported" (usable from other packages). `fmt.println` would not compile.
	fmt.Println("Hello, Go!")

	// runtime.GOOS / GOARCH are constants filled in at COMPILE time with the
	// operating system and CPU architecture this binary was built for,
	// e.g. "darwin arm64" (Apple Silicon Mac) or "linux amd64" (most servers).
	fmt.Println("OS/Arch:", runtime.GOOS, runtime.GOARCH)

	// The version of the Go toolchain that compiled this binary.
	fmt.Println("Go version:", runtime.Version())

	// Number of logical CPU cores available. The scheduler (lesson 19) uses this to
	// decide how many goroutines can truly run at the same instant.
	fmt.Println("CPUs:", runtime.NumCPU())
}

// WHAT HAPPENS WHEN YOU RUN `go build`
//
//  1. PARSE. The compiler reads your .go files and turns the text into a tree
//     structure (an AST) that represents the program.
//  2. TYPE-CHECK. It verifies every expression has a valid type: you cannot add a
//     string to an int, call a function with the wrong arguments, or leave a
//     variable/import unused. Any mistake stops the build HERE, before the program
//     ever runs. This is why Go catches whole classes of bugs that dynamic languages
//     only reveal at runtime.
//  3. COMPILE. Each package is translated into MACHINE CODE: the raw instructions
//     your CPU executes. (Python/JS are usually interpreted, or JIT-compiled while
//     running; Go does all this work ahead of time.)
//  4. LINK. The linker stitches all compiled packages, plus the Go RUNTIME, into a
//     single executable file. The runtime is a chunk of Go code that ships inside
//     EVERY Go program and provides the garbage collector (lesson 18) and the
//     goroutine scheduler (lesson 19).
//
// The result is ONE STATIC BINARY: you can copy that file to another machine with the
// same OS and CPU architecture and run it. No Go installation, interpreter, VM or
// node_modules folder is needed. That is why Go is popular for command-line tools and
// servers (and for tiny Docker images).
//
// Compare with JavaScript: there is no interpreter start-up or JIT warm-up phase;
// the code is already machine code when the program launches.
//
// USEFUL COMMANDS
//   go run .            compile to a temporary binary, run it, throw the binary away
//   go build            compile and keep the binary in the current folder
//   go vet ./...        static checks for code that compiles but looks suspicious
//                       (e.g. a Printf whose format does not match its arguments)
//   gofmt -w .          rewrite files into the ONE official code style, so teams
//                       never argue about formatting
//   GOOS=linux GOARCH=amd64 go build
//                       CROSS-COMPILE: build a Linux binary while sitting on a Mac.
//                       Go's compiler can target any supported OS/CPU with two
//                       environment variables, no extra toolchain needed.
