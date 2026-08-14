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
// TODO (Exercise 5): as written, this function reproduces the classic
// nil-pointer-in-an-interface gotcha. Run it once to see err == nil
// print false even on the success path (currentlyHeld < maxBooksAllowed).
// Then fix CheckOut so the success path returns a literal nil instead of
// a *CheckoutError variable that happens to be nil - don't let the
// concrete pointer type leak into the error return slot.
func CheckOut(book string, currentlyHeld, maxBooksAllowed int) error {
	var problem *CheckoutError = nil
	if currentlyHeld >= maxBooksAllowed {
		problem = &CheckoutError{Book: book}
	}
	return problem
}
