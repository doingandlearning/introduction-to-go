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

// Volunteer embeds Person - no field name, which is what triggers
// promotion. Volunteer.Name and (absent the override below)
// Volunteer.Greet() would both resolve to Person's, automatically.
type Volunteer struct {
	Person
	ShiftHours int
}

// Greet shadows the promoted Person.Greet. This is not a polymorphic
// override - there's no virtual dispatch involved, no shared base-class
// pointer being resolved through. It's a distinct method that happens to
// share a name, and it calls through to the embedded Person's Greet()
// explicitly - the manual "super call."
func (v Volunteer) Greet() string {
	return v.Person.Greet() + " I'm volunteering today."
}
