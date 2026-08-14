package catalog

// ResetCopies takes a Book BY VALUE, sets Copies to 0 on its local copy,
// and returns the copy. The caller's original Book is untouched.
//
// TODO(Exercise 3): implement — set b.Copies = 0 and return b.
func ResetCopies(b Book) Book {
	// TODO: set b.Copies = 0, then return b.
	return b
}

// ResetCopiesPtr takes a *Book and sets Copies to 0 THROUGH the pointer,
// mutating the caller's original Book.
//
// TODO(Exercise 3): implement — set b.Copies = 0 through the pointer.
func ResetCopiesPtr(b *Book) {
	// TODO: set b.Copies = 0 through the pointer (b.Copies = 0 works
	// directly — Go dereferences struct pointers automatically for
	// field access).
}
