package catalog

import "testing"

func TestNewBook_Success(t *testing.T) {
	b, err := NewBook("The Go Programming Language", "Donovan & Kernighan", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Title != "The Go Programming Language" {
		t.Errorf("Title = %q, want %q", b.Title, "The Go Programming Language")
	}
	if b.Author != "Donovan & Kernighan" {
		t.Errorf("Author = %q, want %q", b.Author, "Donovan & Kernighan")
	}
	if b.CopiesAvailable != 3 {
		t.Errorf("CopiesAvailable = %d, want 3", b.CopiesAvailable)
	}
	if b.TimesBorrowed != 0 {
		t.Errorf("TimesBorrowed = %d, want 0", b.TimesBorrowed)
	}
}

func TestNewBook_RejectsNegativeCopies(t *testing.T) {
	b, err := NewBook("Bad Book", "Nobody", -1)
	if err == nil {
		t.Fatal("expected an error for negative copies, got nil")
	}
	if b != nil {
		t.Errorf("expected a nil *Book on error, got %+v", b)
	}
}

func TestEstimatedReadHours_ValueReceiverDoesNotMutate(t *testing.T) {
	b := Book{Title: "Some Book", PageCount: 400}

	got := b.EstimatedReadHours()

	if got != 10 {
		t.Errorf("EstimatedReadHours() = %v, want 10", got)
	}
	if b.PageCount != 400 {
		t.Errorf("PageCount = %d after calling a value-receiver method, want unchanged 400", b.PageCount)
	}
}

func TestCheckout_PointerReceiverMutates(t *testing.T) {
	b := Book{Title: "Some Book", CopiesAvailable: 2, TimesBorrowed: 0}

	err := b.Checkout()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.CopiesAvailable != 1 {
		t.Errorf("CopiesAvailable = %d, want 1 after Checkout", b.CopiesAvailable)
	}
	if b.TimesBorrowed != 1 {
		t.Errorf("TimesBorrowed = %d, want 1 after Checkout", b.TimesBorrowed)
	}
}

func TestCheckout_ErrorWhenNoCopiesAvailable(t *testing.T) {
	b := Book{Title: "Some Book", CopiesAvailable: 0, TimesBorrowed: 0}

	err := b.Checkout()

	if err == nil {
		t.Fatal("expected an error when no copies are available, got nil")
	}
	if b.CopiesAvailable != 0 {
		t.Errorf("CopiesAvailable = %d, want unchanged 0", b.CopiesAvailable)
	}
	if b.TimesBorrowed != 0 {
		t.Errorf("TimesBorrowed = %d, want unchanged 0", b.TimesBorrowed)
	}
}
