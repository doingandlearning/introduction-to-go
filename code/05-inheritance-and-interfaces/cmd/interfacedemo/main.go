// Command interfacedemo shows interfaces satisfied implicitly by three
// unrelated types - no shared base type, no "implements" keyword,
// no declared link between the types and the interface at all.
//
//	go run ./cmd/interfacedemo
package main

import "fmt"

// Greeter is a one-method interface. Idiomatic Go interfaces are small,
// often exactly like this - compare io.Reader, io.Writer, fmt.Stringer,
// each one method in the standard library.
type Greeter interface {
	Greet() string
}

// Person satisfies Greeter.
type Person struct {
	Name string
}

func (p Person) Greet() string {
	return "Hi, I'm " + p.Name
}

// Robot satisfies Greeter too. Robot's author never saw the Greeter
// interface declared anywhere - it doesn't matter. The method set
// matches, so it qualifies. No relationship to Person whatsoever.
type Robot struct {
	ID string
}

func (r Robot) Greet() string {
	return "BEEP BOOP HELLO, UNIT " + r.ID
}

// Parrot also satisfies Greeter, independently again.
type Parrot struct {
	Vocabulary string
}

func (p Parrot) Greet() string {
	return "Squawk! " + p.Vocabulary
}

// WelcomeAll takes anything satisfying Greeter and calls Greet on it.
// It has no idea, and no need to know, what concrete types it's holding.
func WelcomeAll(greeters []Greeter) {
	for _, g := range greeters {
		fmt.Println(g.Greet())
	}
}

func main() {
	greeters := []Greeter{
		Person{Name: "Ada"},
		Robot{ID: "7"},
		Parrot{Vocabulary: "Polly wants a byte slice"},
	}

	WelcomeAll(greeters)

	// describe demonstrates the any + type switch escape hatch for when
	// even a one-method interface is too narrow.
	describe(42)
	describe("hello")
	describe(3.14)
}

func describe(x any) {
	switch v := x.(type) {
	case int:
		fmt.Println("int:", v*2)
	case string:
		fmt.Println("string:", v+"!")
	default:
		fmt.Printf("something else: %v\n", v)
	}
}
