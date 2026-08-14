// Package catalog models books in a small library system for Lab 4.
package catalog

// Book is a single title held by the library.
type Book struct {
	Title           string
	Author          string
	PageCount       int
	CopiesAvailable int
	TimesBorrowed   int
}

// EstimatedReadHours returns a rough read-time estimate based on page
// count, at an assumed 40 pages per hour.
//
// TODO (Exercise 1): implement this. It should return
// float64(b.PageCount) / 40. This is a VALUE receiver on purpose — it
// only computes from b, it never needs to change it.
func (b Book) EstimatedReadHours() float64 {
	// TODO: replace this placeholder.
	return 0
}

// Checkout records a loan: it decrements CopiesAvailable and increments
// TimesBorrowed. It must use a POINTER receiver, because it mutates the
// Book — see Exercise 2.
//
// TODO (Exercise 2): implement this. It should:
//  1. Return an error if CopiesAvailable is already 0 (nothing to
//     lend) — use fmt.Errorf and mention the book's Title.
//  2. Otherwise decrement CopiesAvailable by 1, increment
//     TimesBorrowed by 1, and return nil.
func (b *Book) Checkout() error {
	// TODO: replace this placeholder.
	return nil
}

// NewBook builds a Book, rejecting an invalid initial copy count.
//
// TODO (Exercise 3): implement this. It should:
//  1. Return an error if copies is negative (use fmt.Errorf, and
//     include the rejected value in the message).
//  2. Otherwise return a *Book with the given fields and
//     TimesBorrowed: 0.
func NewBook(title, author string, copies int) (*Book, error) {
	// TODO: replace this placeholder. It returns a non-nil zero-valued
	// Book for now so callers don't crash on a nil dereference before
	// you've implemented validation — that's not the real behavior.
	return &Book{}, nil
}

// Return records a returned copy: it increments CopiesAvailable.
//
// Exercise 5: this is deliberately written with a VALUE receiver, even
// though Checkout above uses a pointer receiver on the same type. Run
// it, confirm the mutation doesn't stick, then fix the receiver.
func (b Book) Return() {
	b.CopiesAvailable++
}
