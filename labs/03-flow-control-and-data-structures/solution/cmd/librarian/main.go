// Command librarian is the entry point for Lab 3: a small library
// circulation desk exercising loops, switches, struct semantics, slice
// aliasing, and maps.
//
//	go run ./cmd/librarian
package main

import (
	"fmt"
	"strings"

	"example.com/librarian/internal/catalog"
)

func main() {
	books := catalog.SampleCatalog()

	// -- Exercise 1: loops --
	fmt.Println("-- catalog --")
	catalog.PrintCatalog(books)

	fmt.Println("available:", catalog.CountAvailable(books))

	if first, ok := catalog.FindFirstAvailable(books); ok {
		fmt.Println("first available:", first.Title)
	} else {
		fmt.Println("nothing available")
	}

	// -- Exercise 2: switches --
	fmt.Println()
	fmt.Println("-- late fees --")
	for _, days := range []int{2, 10, 45} {
		fmt.Printf("%d days late -> %s\n", days, catalog.LateFeeTier(days))
	}

	fmt.Println()
	fmt.Println("-- desk schedule --")
	for day := 1; day <= 7; day++ {
		fmt.Printf("day %d -> %s\n", day, catalog.DeskSchedule(day))
	}

	// -- Exercise 3: struct value vs. pointer --
	fmt.Println()
	fmt.Println("-- reset copies: value vs pointer --")
	book := books[0]
	fmt.Println("before:", book.Copies)

	_ = catalog.ResetCopies(book)
	fmt.Println("after ResetCopies (by value):", book.Copies) // unchanged

	catalog.ResetCopiesPtr(&book)
	fmt.Println("after ResetCopiesPtr (by pointer):", book.Copies) // now 0

	// -- Exercise 4: reproduce the slice-aliasing bug --
	fmt.Println()
	fmt.Println("-- waitlist aliasing --")

	waitlist := []string{"Alice", "Bob", "Carla", "Dev", "Eve"}
	fmt.Println("waitlist:", waitlist, "len", len(waitlist), "cap", cap(waitlist))

	nextUp := waitlist[1:3] // ["Bob", "Carla"]
	fmt.Println("nextUp:  ", nextUp, "len", len(nextUp), "cap", cap(nextUp))

	nextUp[0] = "REPLACED"
	fmt.Println("after nextUp[0] = \"REPLACED\":")
	fmt.Println("  waitlist:", waitlist, "<- changed through nextUp")

	// Still within capacity: this append overwrites waitlist[3] silently.
	nextUp = append(nextUp, "Frank")
	fmt.Println("after append(nextUp, \"Frank\") (within capacity):")
	fmt.Println("  nextUp:  ", nextUp, "len", len(nextUp), "cap", cap(nextUp))
	fmt.Println("  waitlist:", waitlist, "<- waitlist[3] silently overwritten")

	capBefore := cap(nextUp)
	nextUp = append(nextUp, "George", "Hana", "Ivan") // exceeds capacity now
	capAfter := cap(nextUp)
	fmt.Printf("cap before exceeding: %d, cap after: %d (reallocated: %v)\n",
		capBefore, capAfter, capAfter > capBefore)

	nextUp[0] = "NO LONGER SHARED"
	fmt.Println("after nextUp[0] = \"NO LONGER SHARED\":")
	fmt.Println("  waitlist:", waitlist, "<- unaffected, sharing broke at the reallocation")

	// -- Exercise 5: word-frequency counter with comma-ok --
	fmt.Println()
	fmt.Println("-- checkout log --")

	log := "pratchett tolkien pratchett le-guin pratchett tolkien"
	counts := make(map[string]int)
	for _, author := range strings.Fields(log) {
		counts[author]++
	}
	counts["asimov"] = 0 // explicitly tracked, zero occurrences

	for _, name := range []string{"pratchett", "vonnegut", "asimov"} {
		n, ok := counts[name]
		switch {
		case !ok:
			fmt.Printf("%-10s never checked out\n", name)
		default:
			fmt.Printf("%-10s %d checkout(s)\n", name, n)
		}
	}
	// "vonnegut" was never written to the map -> ok is false.
	// "asimov" was explicitly set to 0 -> ok is true, n is 0.
	// Both print n == 0-ish results, but only comma-ok tells them apart.

	// -- Exercise 6: map iteration order --
	fmt.Println()
	fmt.Println("-- map order --")

	fmt.Print("pass 1: ")
	for k := range counts {
		fmt.Print(k, " ")
	}
	fmt.Println()

	fmt.Print("pass 2: ")
	for k := range counts {
		fmt.Print(k, " ")
	}
	fmt.Println()
	fmt.Println("Run this program again as a separate execution and compare — order is randomized on purpose.")
}
