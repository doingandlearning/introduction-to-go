// This test file is already complete — you're not writing it. It's the
// specification for Exercises 1 and 2: run `go test ./...` now, before
// touching library.go, and it won't even compile yet. Implement the
// embedding (Exercise 1) and the shadowed Greet() (Exercise 2) until it
// passes. Writing a test like this yourself is Topic 12's job, not
// this one's.
package library

import "testing"

// TestVolunteerPromotion confirms that embedding Person in Volunteer
// promotes Name onto Volunteer, and that Volunteer's own Greet() shadows
// the promoted one by calling through to it explicitly (Exercises 1-2).
func TestVolunteerPromotion(t *testing.T) {
	v := Volunteer{
		Person:     Person{Name: "Sam"},
		ShiftHours: 4,
	}

	if v.Name != "Sam" {
		t.Errorf("v.Name = %q, want %q (Person.Name should be promoted)", v.Name, "Sam")
	}

	promoted := v.Person.Greet()
	shadowed := v.Greet()

	if shadowed == promoted {
		t.Errorf("v.Greet() should shadow v.Person.Greet(), got identical strings: %q", shadowed)
	}

	want := promoted + " I'm volunteering today."
	if shadowed != want {
		t.Errorf("v.Greet() = %q, want %q (should call through to v.Person.Greet())", shadowed, want)
	}
}
