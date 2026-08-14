// Command library exercises the internal/catalog package for Lab 4:
// value vs. pointer receivers, the map-of-structs gotcha, a validating
// constructor, the zero value, and a mixed-receiver bug found and fixed.
//
//	go run ./cmd/library
package main

import (
	"fmt"

	"example.com/library/internal/catalog"
)

func main() {
	// --- Exercise 1: value-receiver method -----------------------------
	fmt.Println("--- Exercise 1 ---")
	b1 := catalog.Book{
		Title:     "The Go Programming Language",
		Author:    "Donovan & Kernighan",
		PageCount: 380,
	}
	fmt.Printf("%q is about %.1f hours of reading\n", b1.Title, b1.EstimatedReadHours())

	// --- Exercise 2: pointer-receiver method + the map gotcha ----------
	fmt.Println("--- Exercise 2 ---")
	b2 := catalog.Book{Title: "Go in Action", CopiesAvailable: 2}
	fmt.Printf("before checkout: %d copies available\n", b2.CopiesAvailable)
	if err := b2.Checkout(); err != nil {
		fmt.Println("checkout failed:", err)
	}
	fmt.Printf("after checkout: %d copies available (mutation stuck)\n", b2.CopiesAvailable)

	library := map[string]catalog.Book{
		"go-in-action": {Title: "Go in Action", CopiesAvailable: 2},
		"go-web-exam":  {Title: "Go Web Examples", CopiesAvailable: 1},
	}

	// DOES NOT COMPILE — uncomment to trigger it live:
	// library["go-in-action"].Checkout()
	//
	// error: cannot call pointer method on library["go-in-action"]
	//        cannot take address of library["go-in-action"]
	//
	// A map value has no stable address (the map can rehash and move
	// it), so Go refuses to let a pointer-receiver method run directly
	// on one. Fix: read the value out, mutate the local copy, write it
	// back under the same key.
	entry := library["go-in-action"]
	if err := entry.Checkout(); err != nil {
		fmt.Println("checkout failed:", err)
	}
	library["go-in-action"] = entry
	fmt.Printf("go-in-action copies available after fix: %d\n", library["go-in-action"].CopiesAvailable)

	// --- Exercise 3: validating constructor -----------------------------
	fmt.Println("--- Exercise 3 ---")
	good, err := catalog.NewBook("Clean Code", "Robert C. Martin", 3)
	if err != nil {
		fmt.Println("unexpected error:", err)
	} else {
		fmt.Printf("created: %+v\n", *good)
	}

	bad, err := catalog.NewBook("Broken Book", "Nobody", -1)
	if err != nil {
		fmt.Println("rejected as expected:", err)
	} else {
		fmt.Printf("should not have succeeded: %+v\n", *bad)
	}

	// --- Exercise 4: the zero value is a real value ---------------------
	fmt.Println("--- Exercise 4 ---")
	var zero catalog.Book
	fmt.Printf("zero-value book: %+v\n", zero)
	fmt.Printf("its EstimatedReadHours(): %.1f (runs fine, no crash)\n", zero.EstimatedReadHours())

	// --- Exercise 5: mixed receivers, on purpose -------------------------
	fmt.Println("--- Exercise 5 ---")
	b5 := catalog.Book{Title: "Returned Copy", CopiesAvailable: 1}
	fmt.Printf("before Return(): %d copies available\n", b5.CopiesAvailable)
	b5.Return()
	fmt.Printf("after Return(): %d copies available\n", b5.CopiesAvailable)
	fmt.Println("Return now uses a pointer receiver, matching Checkout — the increment sticks.")
}
