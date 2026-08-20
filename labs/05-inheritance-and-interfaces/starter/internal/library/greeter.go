package library

// Greeter is satisfied by anything that can greet a visitor - a human,
// a machine, a mascot in a costume, whatever. No shared base type
// required, no declared relationship needed.
//
// TODO (Exercise 3): define this interface with one method,
// Greet() string.
type Greeter interface {
	// TODO: add the Greet() string method signature
}

// Kiosk is a self-service checkout machine.
//
// TODO (Exercise 3): give Kiosk a Greet() method (pointer or value
// receiver, your choice) that returns something like
// "SCAN YOUR CARD TO BEGIN." Kiosk should satisfy Greeter with zero
// declared relationship to it.
type Kiosk struct {
	StationNumber int
}

// Mascot is the library's reading-time costume mascot.
//
// TODO (Exercise 3): give Mascot a Greet() method that returns
// something like "*waves enthusiastically, says nothing*". Mascot
// should satisfy Greeter independently of Person and Kiosk.
type Mascot struct {
	CharacterName string
}

// WelcomeAll prints a greeting for every Greeter in the slice. It has no
// idea what concrete types it's holding, and doesn't need to.
//
// TODO (Exercise 3): uncomment this once Greeter has a Greet() method
// above — as written it won't compile yet, since g.Greet() isn't part
// of Greeter's method set until you add it. You'll need to add the
// "fmt" import back to the top of this file too.
//
// func WelcomeAll(greeters []Greeter) {
// 	for _, g := range greeters {
// 		fmt.Println(g.Greet())
// 	}
// }
