package library

import "fmt"

// Greeter is satisfied by anything that can greet a visitor - a human,
// a machine, a mascot in a costume, whatever. No shared base type
// required, no declared relationship needed. Person (above), Kiosk, and
// Mascot (below) each satisfy it independently.
type Greeter interface {
	Greet() string
}

// Kiosk is a self-service checkout machine. It satisfies Greeter with
// zero declared relationship to the interface or to Person.
type Kiosk struct {
	StationNumber int
}

func (k Kiosk) Greet() string {
	return fmt.Sprintf("SCAN YOUR CARD AT STATION %d TO BEGIN.", k.StationNumber)
}

// Mascot is the library's reading-time costume mascot. It satisfies
// Greeter too, independently of both Person and Kiosk.
type Mascot struct {
	CharacterName string
}

func (m Mascot) Greet() string {
	return "*" + m.CharacterName + " waves enthusiastically, says nothing*"
}

// WelcomeAll prints a greeting for every Greeter in the slice. It has no
// idea what concrete types it's holding, and doesn't need to.
func WelcomeAll(greeters []Greeter) {
	for _, g := range greeters {
		fmt.Println(g.Greet())
	}
}
