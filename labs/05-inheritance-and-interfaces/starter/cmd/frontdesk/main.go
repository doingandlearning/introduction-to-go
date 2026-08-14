// Command frontdesk drives the Lab 5 exercises: struct embedding with
// promotion, an interface satisfied implicitly by unrelated types, an
// any + type switch, and the nil-pointer-in-an-interface gotcha.
//
//	go run ./cmd/frontdesk
package main

import (
	"fmt"

	"example.com/library-frontdesk/internal/library"
)

func main() {
	// --- Exercise 1 & 2: embedding and promotion ---
	p := library.Person{Name: "Ada"}
	fmt.Println(p.Greet())

	// Uncomment once Volunteer embeds Person (Exercise 1):
	// v := library.Volunteer{Person: library.Person{Name: "Sam"}, ShiftHours: 4}
	// fmt.Println("Volunteer.Name:  ", v.Name)
	// fmt.Println("Volunteer.Greet():", v.Greet())

	// --- Exercise 3: Greeter interface, three unrelated types ---
	// Uncomment once Greeter, Kiosk, and Mascot are implemented:
	// greeters := []library.Greeter{
	// 	library.Person{Name: "Ada"},
	// 	library.Kiosk{StationNumber: 2},
	// 	library.Mascot{CharacterName: "Bookworm Barry"},
	// }
	// library.WelcomeAll(greeters)

	// --- Exercise 4: any + type switch ---
	logCheckIn(42)
	logCheckIn("Priya")
	logCheckIn(3.14)

	// --- Exercise 5: the nil-interface gotcha ---
	err := library.CheckOut("The Go Programming Language", 0, 5)
	fmt.Println("CheckOut() err == nil ->", err == nil)
}

// logCheckIn accepts any front-desk log event and reports what kind it
// got. int is treated as a visitor count, string as a patron name.
//
// TODO (Exercise 4): implement this with a type switch
// (switch v := x.(type)) handling at least int, string, and a default
// case.
func logCheckIn(x any) {
	// TODO: replace this placeholder with a real type switch.
	fmt.Println("check-in event:", x)
}
