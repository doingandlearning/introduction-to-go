// Command library exercises the internal/catalog package for Lab 4:
// value vs. pointer receivers, the map-of-structs gotcha, a validating
// constructor, the zero value, and a mixed-receiver bug to find and fix.
//
// This file compiles and runs as-is, even before you've implemented
// anything in internal/catalog/catalog.go — the stub methods there
// return zero values, so the output will look wrong until you fill
// them in one exercise at a time. Re-run `go run ./cmd/library` after
// each exercise to see the output change.
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
	fmt.Printf("after checkout: %d copies available\n", b2.CopiesAvailable)

	library := map[string]catalog.Book{
		"go-in-action": {Title: "Go in Action", CopiesAvailable: 2},
		"go-web-exam":  {Title: "Go Web Examples", CopiesAvailable: 1},
	}

	// TODO: uncomment the next line and run `go build ./cmd/library`.
	// Read the compiler error, then re-comment this line — it will
	// never compile, that's the point.
	// library["go-in-action"].Checkout()

	// TODO: fix it with read-mutate-write-back: read the Book out of
	// the map into a local variable, call Checkout on the local
	// variable, then write it back into the map under the same key.
	// Print library["go-in-action"].CopiesAvailable afterward to
	// confirm it changed.

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
}
