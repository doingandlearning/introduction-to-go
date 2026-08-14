package catalog

import "testing"

func TestResetCopies(t *testing.T) {
	original := Book{Title: "The Go Programming Language", Copies: 2}

	reset := ResetCopies(original)

	if reset.Copies != 0 {
		t.Errorf("ResetCopies(...).Copies = %d, want 0", reset.Copies)
	}
	if original.Copies != 2 {
		t.Errorf("original.Copies = %d, want unchanged 2 — ResetCopies must not mutate its argument", original.Copies)
	}
}

func TestResetCopiesPtr(t *testing.T) {
	book := Book{Title: "Design Patterns", Copies: 5}

	ResetCopiesPtr(&book)

	if book.Copies != 0 {
		t.Errorf("book.Copies = %d, want 0 after ResetCopiesPtr", book.Copies)
	}
}
