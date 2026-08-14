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
