// Lesson 11: Composition over inheritance, using struct and interface embedding.
//
// WHAT YOU WILL LEARN
//   - What embedding is and how method "promotion" works
//   - How to override a promoted method and still reach the original
//   - Interface embedding, for building big interfaces from small ones
//   - Why embedding is NOT inheritance
//
// BACKGROUND: in languages with inheritance you write "a Car IS-A Vehicle". Go
// prefers COMPOSITION: "a Car HAS an Engine and a Logger". Embedding is syntax sugar
// that makes composition convenient, so you do not have to write car.Engine.Start()
// everywhere.
package main

import "fmt"

type Logger struct{ Prefix string }

func (l Logger) Log(msg string) { fmt.Println(l.Prefix + msg) }

type Engine struct{ HP int }

func (e Engine) Start() string { return fmt.Sprintf("engine with %dhp running", e.HP) }

// EMBEDDING: a field written with only a type and NO name is an embedded field.
// Car stores a real Logger and a real Engine inside it (memory layout: Logger, then
// Engine, then Name). The compiler then PROMOTES the embedded types' fields and
// methods, so `car.Log("x")` is automatically translated to `car.Logger.Log("x")`
// and `car.HP` to `car.Engine.HP`. The embedded value is still reachable by its type
// name (c.Engine).
type Car struct {
	Logger // embedded: no field name
	Engine
	Name string
}

// OVERRIDING: if the outer type declares a method with the same name, the OUTER one
// wins when you call c.Start(). The inner one is not gone: c.Engine.Start() still
// reaches it (used here to build on top of it). The mechanism is simply "the shallowest
// name wins"; nothing is replaced.
func (c Car) Start() string { return c.Name + ": " + c.Engine.Start() }

// INTERFACE EMBEDDING: listing Reader and Writer inside ReadWriter means "ReadWriter
// requires all the methods of Reader AND all the methods of Writer." This is how the
// standard library builds io.ReadWriter, io.ReadCloser, and others from small parts.
type Reader interface{ Read() string }
type Writer interface{ Write(string) }
type ReadWriter interface {
	Reader
	Writer
}

type File struct{ data string }

// Pointer receivers because Write changes the struct; hence only *File (not File)
// satisfies the interfaces (lesson 10, method sets).
func (f *File) Read() string   { return f.data }
func (f *File) Write(s string) { f.data += s }

func main() {
	// Positional struct literal: values fill the fields in order: Logger, Engine, Name.
	// (Fine for a tiny demo; prefer named fields in real code.)
	c := Car{Logger{"[car] "}, Engine{300}, "Mustang"}
	c.Log("created")       // promoted from Logger: prints "[car] created"
	fmt.Println(c.HP)      // promoted field: 300
	fmt.Println(c.Start()) // Car.Start wins over Engine.Start

	// *File has both Read and Write, so it satisfies ReadWriter implicitly.
	var rw ReadWriter = &File{}
	rw.Write("hi ")
	rw.Write("there")
	fmt.Println(rw.Read()) // hi there

	// IMPORTANT: this is NOT inheritance. A Car is not an Engine and cannot be passed
	// where an Engine is expected. There is also no "virtual" dispatch back to the
	// outer type: when Engine's own methods run, their receiver is the Engine, and
	// they know nothing about Car. If Engine.Start called a method named Describe, it
	// would always get Engine's Describe, even if Car defined one. Want polymorphism?
	// Use interfaces.
}
