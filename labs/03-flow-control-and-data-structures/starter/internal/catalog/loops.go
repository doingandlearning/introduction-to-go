package catalog

import "io"

// PrintCatalog prints each book's title and copy count to w, using a
// classic C-style for loop.
//
// TODO(Exercise 1): implement this. It should:
//  1. Add an `import "fmt"` line above `package catalog` — this file
//     doesn't import it yet, since nothing here needs it until you do.
//  2. Loop over books with a classic C-style for loop —
//     `for i := 0; i < len(books); i++ { ... }` — printing each one to w
//     with fmt.Fprintf, something like: "  <title> — <copies> cop(y/ies)".
func PrintCatalog(w io.Writer, books []Book) {
	// TODO: replace this placeholder.
	_ = books
}

// CountAvailable returns how many books have Copies > 0.
//
// TODO(Exercise 1): implement using a while-style for loop —
// a bare condition, no init or post statement.
func CountAvailable(books []Book) int {
	// TODO: replace this with a while-style for loop that counts books
	// with Copies > 0.
	_ = books
	return 0
}

// FindFirstAvailable scans books for the first one with Copies > 0.
// Returns the zero Book and false if none is found.
//
// TODO(Exercise 1): implement using an infinite for {} loop with an
// explicit break when found (or when you've scanned everything).
func FindFirstAvailable(books []Book) (Book, bool) {
	// TODO: replace this with an infinite for {} loop, indexing manually
	// and breaking when you find a book with Copies > 0, or when you run
	// out of books to check.
	_ = books
	return Book{}, false
}
