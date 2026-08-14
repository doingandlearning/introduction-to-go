// Package catalog models a small library circulation desk: books on the
// shelf, and the operations Lab 3 builds around them.
package catalog

// Book is a single catalog entry. Deliberately plain — a typed bag of
// fields, no methods (those come in a later topic).
type Book struct {
	Title  string
	Copies int
}

// SampleCatalog returns a small fixed catalog used by the lab's demo
// program.
func SampleCatalog() []Book {
	return []Book{
		{Title: "The Go Programming Language", Copies: 2},
		{Title: "Design Patterns", Copies: 0},
		{Title: "The Pragmatic Programmer", Copies: 1},
		{Title: "Clean Code", Copies: 0},
		{Title: "Introduction to Algorithms", Copies: 3},
	}
}
