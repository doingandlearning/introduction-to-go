// Package library models the front-desk greeting system for Sunnyvale
// Public Library: staff, volunteers, self-service kiosks, and a reading-
// time mascot all need to be able to greet a visitor.
package library

import "fmt"

// Person is the base type for anyone who works the front desk.
type Person struct {
	Name string
}

// Greet returns a standard front-desk greeting.
func (p Person) Greet() string {
	return fmt.Sprintf("Hi, I'm %s, welcome to the library!", p.Name)
}

// Volunteer embeds Person.
//
// TODO (Exercise 1): embed Person here with no field name, so
// Volunteer.Name and Volunteer.Greet() are promoted automatically.
// Add a ShiftHours int field alongside it.
type Volunteer struct {
	// TODO: embed Person
	ShiftHours int
}

// TODO (Exercise 2): give Volunteer its own Greet() method that shadows
// the promoted one. It should call through to the embedded Person's
// Greet() explicitly and append ", volunteering today" to the result.
