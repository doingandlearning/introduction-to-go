// Command switchdemo shows Go's conditionless switch (grading scores) and
// then the fallthrough gotcha: cases do NOT fall through by default, the
// opposite of C/Java/JavaScript, so falling through is an explicit,
// deliberate choice.
//
//	go run ./cmd/switchdemo
package main

import "fmt"

// grade turns a numeric score into a letter grade using a switch with no
// condition at all — it reads like a clean if/else chain, but is a switch.
func grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	default:
		return "F"
	}
}

// openHours reports whether the library is open on a given day, using
// deliberate fallthrough to group Saturday's "open, but only for the
// morning" case into Sunday's "closed" case as a demonstration — this is
// a contrived use, real code reaches for multiple values per case instead
// (case 1, 2, 3, 4, 5:) unless it genuinely wants cascading behavior.
func openHours(day int) string {
	switch day {
	case 1, 2, 3, 4, 5: // Monday-Friday
		return "open 9am-6pm"
	case 6: // Saturday
		// Without fallthrough, this case would end here. With it, we
		// deliberately also run Sunday's case body.
		fallthrough
	case 7: // Sunday
		return "open 10am-2pm only"
	default:
		return "unknown day"
	}
}

func main() {
	fmt.Println("-- grading (conditionless switch) --")
	for _, score := range []int{95, 82, 71, 40} {
		fmt.Printf("score %3d -> grade %s\n", score, grade(score))
	}

	fmt.Println()
	fmt.Println("-- day-of-week (fallthrough) --")
	for day := 1; day <= 7; day++ {
		fmt.Printf("day %d -> %s\n", day, openHours(day))
	}

	// Watch what happens to Saturday (day 6) specifically: it hits case 6,
	// then falls through into case 7's body and returns THAT string, not
	// its own. Comment out the `fallthrough` line above and re-run to see
	// day 6 fall back to returning nothing meaningful from its own case —
	// try it and see the difference for yourself.
}
