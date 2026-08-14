// Command embeddemo shows struct embedding: fields and methods promoted
// from an embedded type, and what happens when the outer type shadows
// one of those methods instead of inheriting it.
//
//	go run ./cmd/embeddemo
package main

import "fmt"

// Person has a name and a method that greets using it.
type Person struct {
	Name string
}

func (p Person) Greet() string {
	return "Hi, I'm " + p.Name
}

// Employee embeds Person with no field name. That's what triggers
// promotion: Employee.Name and Employee.Greet() both become valid,
// forwarding to the embedded Person.
type Employee struct {
	Person
	Role string
}

// Greet shadows the promoted Person.Greet. This is not an override in
// the polymorphic sense - there's no virtual dispatch happening. It's a
// new method on Employee that happens to have the same name, and it
// calls through to the embedded version explicitly, the way super.Greet()
// would in Java, except it's just a normal field access and method call.
func (e Employee) Greet() string {
	return e.Person.Greet() + ", I work here"
}

func main() {
	p := Person{Name: "Ada"}
	fmt.Println("Person.Greet():       ", p.Greet())

	e := Employee{Person: Person{Name: "Sam"}, Role: "Engineer"}
	fmt.Println("Employee.Name:        ", e.Name)          // promoted field
	fmt.Println("Employee.Greet():     ", e.Greet())        // shadowed method
	fmt.Println("Employee.Person.Greet():", e.Person.Greet()) // the original, reached explicitly
}
