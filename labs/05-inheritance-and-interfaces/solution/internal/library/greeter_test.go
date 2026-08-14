package library

import "testing"

// TestGreeterSatisfiedByMultipleTypes confirms that Person, Kiosk, and
// Mascot each satisfy Greeter independently, with no declared
// relationship between them (Exercise 3).
func TestGreeterSatisfiedByMultipleTypes(t *testing.T) {
	greeters := []Greeter{
		Person{Name: "Ana"},
		Kiosk{StationNumber: 3},
		Mascot{CharacterName: "Sparky"},
	}

	want := []string{
		"Hi, I'm Ana, welcome to the library!",
		"SCAN YOUR CARD AT STATION 3 TO BEGIN.",
		"*Sparky waves enthusiastically, says nothing*",
	}

	for i, g := range greeters {
		got := g.Greet()
		if got != want[i] {
			t.Errorf("greeters[%d].Greet() = %q, want %q", i, got, want[i])
		}
	}
}
