# How Go Works

Small, runnable lessons with detailed comments that explain how Go works behind the scenes: from how a program is compiled and scheduled to the core language, interfaces, generics and concurrency. A companion to `how-js-works`, aimed at beginners and at developers coming from other languages.

## How to use

Each lesson is a folder with a `main.go`. Read the comments, run it, then change things and run it again.

```bash
go run ./05-arrays-and-slices     # run one lesson
go vet ./...                      # check everything
go test -v ./20-testing           # run the tests lesson
go run -race ./15-goroutines      # detect data races
```

Requires Go 1.23 or newer.

## Contents

### Part 1 — Getting started
| # | Lesson | Topic |
|---|--------|-------|
| 01 | `01-hello-and-compilation` | Build vs run, `init`/`main`, static binaries, cross-compiling |
| 02 | `02-variables-and-types` | Declarations, zero values, no implicit conversion, `iota` |
| 03 | `03-control-flow` | `if`, the single `for`, `switch`, type switch, labels, per-iteration loop vars |
| 04 | `04-functions` | Multiple returns, variadics, closures, `defer`, function types |

### Part 2 — Data
| # | Lesson | Topic |
|---|--------|-------|
| 05 | `05-arrays-and-slices` | Slice header internals, capacity, `append` aliasing gotchas |
| 06 | `06-maps` | Maps, comma-ok, random iteration order, sets |
| 07 | `07-strings-and-runes` | UTF-8, bytes vs runes, `strings.Builder`, format verbs |
| 08 | `08-structs-and-methods` | Structs, value vs pointer receivers, methods on any type |
| 09 | `09-pointers` | Pointers, pass-by-value, nil dereference |

### Part 3 — Abstraction
| # | Lesson | Topic |
|---|--------|-------|
| 10 | `10-interfaces` | Implicit interfaces, `io.Reader`, the nil-interface gotcha |
| 11 | `11-embedding-and-composition` | Embedding and promotion (not inheritance) |
| 12 | `12-errors` | Errors as values, wrapping, `errors.Is`/`As`/`Join` |
| 13 | `13-defer-panic-recover` | Cleanup and the panic path |
| 14 | `14-generics` | Type parameters, constraints, generic types |

### Part 4 — Concurrency and the runtime
| # | Lesson | Topic |
|---|--------|-------|
| 15 | `15-goroutines` | Goroutines, `WaitGroup`, cheap threads |
| 16 | `16-channels-and-select` | Channels, pipelines, `select`, worker pools |
| 17 | `17-sync-and-context` | Mutex, atomics, `Once`, `context` cancellation |
| 18 | `18-memory-and-gc` | Stack vs heap, escape analysis, the garbage collector |
| 19 | `19-scheduler` | The G-M-P scheduler, `GOMAXPROCS`, parked goroutines |

### Part 5 — Practice
| # | Lesson | Topic |
|---|--------|-------|
| 20 | `20-testing` | Table-driven tests, benchmarks, fuzzing, coverage |
| 21 | `21-stdlib-tour` | JSON, time, HTTP server/client, `slog`, files |

## Where to go next
[A Tour of Go](https://go.dev/tour), [Effective Go](https://go.dev/doc/effective_go), [Go by Example](https://gobyexample.com), the [Go spec](https://go.dev/ref/spec), and the [Go blog](https://go.dev/blog).
