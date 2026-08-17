package catalog

import (
	"fmt"
	"io"
)

// PrintCatalog prints each book's title and copy count to w, using a
// classic C-style for loop.
func PrintCatalog(w io.Writer, books []Book) {
	for i := 0; i < len(books); i++ {
		fmt.Fprintf(w, "  %s — %d cop(y/ies)\n", books[i].Title, books[i].Copies)
	}
}

// CountAvailable returns how many books have Copies > 0, using a
// while-style for loop (bare condition, no init or post statement).
func CountAvailable(books []Book) int {
	count := 0
	i := 0
	for i < len(books) {
		if books[i].Copies > 0 {
			count++
		}
		i++
	}
	return count
}

// FindFirstAvailable scans books for the first one with Copies > 0,
// using an infinite for {} loop with an explicit break.
func FindFirstAvailable(books []Book) (Book, bool) {
	i := 0
	for {
		if i >= len(books) {
			return Book{}, false
		}
		if books[i].Copies > 0 {
			return books[i], true
		}
		i++
	}
}
