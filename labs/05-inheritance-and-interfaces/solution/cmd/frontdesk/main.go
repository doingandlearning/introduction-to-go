// Command frontdesk drives the Lab 5 exercises: struct embedding with
// promotion, an interface satisfied implicitly by unrelated types, an
// any + type switch, and the nil-pointer-in-an-interface gotcha.
//
//	go run ./cmd/frontdesk
package main

import (
	"fmt"
	"io"
	"os"

	"example.com/library-frontdesk/internal/library"
)

func main() {
	// --- Exercise 1 & 2: embedding and promotion ---
	p := library.Person{Name: "Ada"}
	fmt.Println(p.Greet())

	v := library.Volunteer{Person: library.Person{Name: "Sam"}, ShiftHours: 4}
	fmt.Println("Volunteer.Name:   ", v.Name)    // promoted field
	fmt.Println("Volunteer.Greet():", v.Greet()) // shadowed method
	fmt.Println("Person.Greet() via embedded field:", v.Person.Greet())

	// --- Exercise 3: Greeter interface, three unrelated types ---
	greeters := []library.Greeter{
		library.Person{Name: "Ada"},
		library.Kiosk{StationNumber: 2},
		library.Mascot{CharacterName: "Bookworm Barry"},
	}
	library.WelcomeAll(greeters)

	// --- Exercise 4: any + type switch ---
	logCheckIn(os.Stdout, 42)
	logCheckIn(os.Stdout, "Priya")
	logCheckIn(os.Stdout, 3.14)

	// --- Exercise 5: the nil-interface gotcha, before and after ---
	buggyErr := library.CheckOutBuggy("The Go Programming Language", 0, 5)
	fmt.Println("CheckOutBuggy() err == nil ->", buggyErr == nil) // false - the gotcha

	fixedErr := library.CheckOut("The Go Programming Language", 0, 5)
	fmt.Println("CheckOut() err == nil ->      ", fixedErr == nil) // true - fixed
}

// logCheckIn accepts any front-desk log event and reports what kind it
// got. int is treated as a visitor count, string as a patron name. It
// writes to w instead of stdout directly - same reason as Topic 2's
// "fmt doesn't just print to the terminal" slide: an io.Writer
// parameter is what makes this testable with a bytes.Buffer instead of
// spawning the real program.
func logCheckIn(w io.Writer, x any) {
	switch v := x.(type) {
	case int:
		fmt.Fprintf(w, "check-in event: %d visitors just arrived\n", v)
	case string:
		fmt.Fprintf(w, "check-in event: patron %q checked in\n", v)
	default:
		fmt.Fprintf(w, "check-in event: unrecognized event type: %v\n", v)
	}
}
