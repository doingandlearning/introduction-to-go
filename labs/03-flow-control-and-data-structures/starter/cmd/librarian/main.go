// Command librarian is the entry point for Lab 3: a small library
// circulation desk exercising loops, switches, struct semantics, slice
// aliasing, and maps.
//
//	go run ./cmd/librarian
package main

import (
	"fmt"
	"os"

	"example.com/librarian/internal/catalog"
)

func main() {
	books := catalog.SampleCatalog()

	// -- Exercise 1: loops --
	fmt.Println("-- catalog --")
	catalog.PrintCatalog(os.Stdout, books)

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
	fmt.Println("after ResetCopies (by value):", book.Copies) // TODO: confirm unchanged

	catalog.ResetCopiesPtr(&book)
	fmt.Println("after ResetCopiesPtr (by pointer):", book.Copies) // TODO: confirm now 0

	// -- Exercise 4: reproduce the slice-aliasing bug --
	fmt.Println()
	fmt.Println("-- waitlist aliasing --")
	// TODO: waitlist := []string{"Alice", "Bob", "Carla", "Dev", "Eve"}
	// TODO: nextUp := waitlist[1:3]
	// TODO: print len/cap of both
	// TODO: nextUp[0] = "REPLACED"; print waitlist, confirm it changed
	// TODO: append one name to nextUp (within capacity); print waitlist,
	//       confirm a name was silently overwritten
	// TODO: print cap(nextUp) before appending enough names to exceed
	//       capacity, then print cap(nextUp) after — confirm it jumped
	// TODO: mutate nextUp[0] once more; print waitlist, confirm it is now
	//       unaffected

	// -- Exercise 5: word-frequency counter with comma-ok --
	fmt.Println()
	fmt.Println("-- checkout log --")
	// TODO: log := "pratchett tolkien pratchett le-guin pratchett tolkien"
	// TODO: split with strings.Fields, build a map[string]int of counts
	// TODO: look up an author who appears several times, one who never
	//       appears, and (optionally) one explicitly set to 0 beforehand —
	//       use comma-ok to print "<name>: <n> checkout(s)" or
	//       "<name>: never checked out" and confirm the zero-vs-never
	//       cases print differently

	// -- Exercise 6: map iteration order --
	fmt.Println()
	fmt.Println("-- map order --")
	// TODO: range over the counts map from Exercise 5 twice in a row,
	//       printing keys each time
	// TODO: run `go run ./cmd/librarian` again as a separate execution
	//       and compare the order across runs
}
