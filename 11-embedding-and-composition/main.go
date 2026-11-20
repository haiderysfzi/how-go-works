// Lesson 11: Composition over inheritance, using struct and interface embedding.
package main

import "fmt"

type Logger struct{ Prefix string }

func (l Logger) Log(msg string) { fmt.Println(l.Prefix + msg) }

type Engine struct{ HP int }

func (e Engine) Start() string { return fmt.Sprintf("engine with %dhp running", e.HP) }

// Embedding: Car "has a" Logger and Engine, and their methods are PROMOTED.
type Car struct {
	Logger // embedded: no field name
	Engine
	Name string
}

// Overriding: define the method on the outer type. The inner one is still reachable.
func (c Car) Start() string { return c.Name + ": " + c.Engine.Start() }

// Interface embedding builds bigger interfaces from small ones.
type Reader interface{ Read() string }
type Writer interface{ Write(string) }
type ReadWriter interface {
	Reader
	Writer
}

type File struct{ data string }

func (f *File) Read() string   { return f.data }
func (f *File) Write(s string) { f.data += s }

func main() {
	c := Car{Logger{"[car] "}, Engine{300}, "Mustang"}
	c.Log("created")       // promoted from Logger
	fmt.Println(c.HP)      // promoted field
	fmt.Println(c.Start()) // Car.Start wins over Engine.Start

	var rw ReadWriter = &File{}
	rw.Write("hi ")
	rw.Write("there")
	fmt.Println(rw.Read())

	// IMPORTANT: this is NOT inheritance. There is no polymorphism through the embedded
	// type: Engine.Start never calls Car.Start. The receiver is always the Engine.
}
