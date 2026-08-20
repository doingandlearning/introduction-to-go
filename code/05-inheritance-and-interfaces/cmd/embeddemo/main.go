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

func Hi(p *Person) {
	p.Greet()
}

// Greet shadows the promoted Person.Greet. This is not an override in
// the polymorphic sense - there's no virtual dispatch happening. It's a
// new method on Employee that happens to have the same name, and it
// calls through to the embedded version explicitly, the way super.Greet()
// would in Java, except it's just a normal field access and method call.
func (e *Employee) Greet() string {
	return e.Person.Greet() + ", I work here"
}

func (p Person) String() string {
	return "Person: " + p.Name
}

func (e *Employee) ChangeRole(newRole string) {
	e.Role = newRole
}

func NewEmployee(name, role string) (*Employee, error) {
	return &Employee{Person: Person{Name: name}, Role: role}, nil
}

func main() {
	p := Person{Name: "Ada"}
	fmt.Println("Person.Greet():       ", p.Greet())

	e := Employee{
		Person: Person{Name: "Sam"},
		Role:   "Engineer"}

	newE, _ := NewEmployee("Kevin", "Trainer")
	fmt.Println("Employee.Name:        ", e.Name)           // promoted field
	fmt.Println("Employee.Person.Greet():     ", e.Greet()) // shadowed method

	fmt.Println("Employee.Person.Greet():", e.Person.Greet()) // the original, reached explicitly
	fmt.Println(newE)

	var greetings []IGreeter = []IGreeter{&p, &e, newE}
	WelcomeAll(greetings)

	fmt.Println(p)
	fmt.Println(e)
	fmt.Println(newE)
}

type IGreeter interface {
	Greet() string
}

type Reader interface {
	Read() string
}

type Writer interface {
	Write(string)
}

type ReadWriter interface {
	Reader
	Writer
}

type Empty interface{}

func PrintEmpty(e Empty) {
	fmt.Println(e)
}

func ReadAll(r Reader) string {
	return r.Read()
}

func WriteAll(w Writer, s string) {
	w.Write(s)
}

func ReadWriteAll(rw ReadWriter, s string) {
	r := rw.Read()
	fmt.Println("Read:", r)
	rw.Write(s)
}

func WelcomeAll(gs []IGreeter) {
	for _, g := range gs {
		fmt.Println(g.Greet())
	}
}

type Printable interface {
	String() string
}
