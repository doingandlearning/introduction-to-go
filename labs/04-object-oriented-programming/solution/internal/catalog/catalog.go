// Package catalog models books in a small library system for Lab 4.
package catalog

import "fmt"

// Book is a single title held by the library.
type Book struct {
	Title           string
	Author          string
	PageCount       int
	CopiesAvailable int
	TimesBorrowed   int
}

// EstimatedReadHours returns a rough read-time estimate based on page
// count, at an assumed 40 pages per hour. VALUE receiver — it only
// computes from b, it never needs to change it.
func (b Book) EstimatedReadHours() float64 {
	return float64(b.PageCount) / 40
}

// Checkout records a loan: it decrements CopiesAvailable and increments
// TimesBorrowed. POINTER receiver — it mutates the Book.
func (b *Book) Checkout() error {
	if b.CopiesAvailable == 0 {
		return fmt.Errorf("no copies of %q available to check out", b.Title)
	}
	b.CopiesAvailable--
	b.TimesBorrowed++
	return nil
}

// NewBook builds a Book, rejecting an invalid initial copy count.
func NewBook(title, author string, copies int) (*Book, error) {
	if copies < 0 {
		return nil, fmt.Errorf("copies cannot be negative, got %d", copies)
	}
	return &Book{
		Title:           title,
		Author:          author,
		CopiesAvailable: copies,
		TimesBorrowed:   0,
	}, nil
}

// Return records a returned copy: it increments CopiesAvailable.
// POINTER receiver, matching Checkout — see Exercise 5. A value
// receiver here would compile fine and silently do nothing, since the
// increment would happen on a copy that's discarded when Return
// returns.
func (b *Book) Return() {
	b.CopiesAvailable++
}
