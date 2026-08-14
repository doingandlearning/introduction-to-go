package catalog

// ResetCopies takes a Book BY VALUE, sets Copies to 0 on its local copy,
// and returns the copy. The caller's original Book is untouched.
func ResetCopies(b Book) Book {
	b.Copies = 0
	return b
}

// ResetCopiesPtr takes a *Book and sets Copies to 0 THROUGH the pointer,
// mutating the caller's original Book.
func ResetCopiesPtr(b *Book) {
	b.Copies = 0
}
