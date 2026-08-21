package library

// CheckoutError reports a problem checking out a book.
type CheckoutError struct {
	Book string
}

func (e *CheckoutError) Error() string {
	return "could not check out: " + e.Book
}

// CheckOut simulates checking a book out. maxBooksAllowed is the
// patron's limit; currentlyHeld is how many they already have.
//
// Fixed version: the success path returns a literal nil directly,
// instead of a *CheckoutError variable that happens to be nil. That
// keeps a concrete pointer type from ever leaking into the interface
// return slot, so err == nil behaves correctly for callers.
func CheckOut(book string, currentlyHeld, maxBooksAllowed int) error {
	if currentlyHeld >= maxBooksAllowed {
		return &CheckoutError{Book: book}
	}
	return nil
}

// CheckOutBuggy is kept for comparison - this is the version students
// see first, that reproduces the nil-pointer-in-an-interface gotcha.
// var problem is declared as a concrete *CheckoutError; even when it's
// left nil, returning it as an error stores the pair (*CheckoutError,
// nil), and err == nil is false because the type half is non-nil.
func CheckOutBuggy(book string, currentlyHeld, maxBooksAllowed int) error {
	if currentlyHeld >= maxBooksAllowed {
		return &CheckoutError{Book: book}
	}
	return nil
}
