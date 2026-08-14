// Exercise 5: prove dependency injection actually buys testability, by
// constructing Service with a fake Repository defined right here in the
// test - no mocking framework, no reflection, no real database.
package catalog

import (
	"errors"
	"testing"
)

// fakeRepository is a hand-written test double for Repository. It exists
// only inside this test file - the same technique cmd/cataloglab uses at
// runtime, but here it's wired up at compile time by go test.
type fakeRepository struct {
	book *Book
	err  error
}

func (f *fakeRepository) FindByISBN(isbn string) (*Book, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.book, nil
}

func TestServiceDescribeSuccess(t *testing.T) {
	fake := &fakeRepository{book: &Book{ISBN: "978-0134190440", Title: "The Go Programming Language"}}
	svc := NewService(fake)

	got, err := svc.Describe("978-0134190440")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "The Go Programming Language (ISBN: 978-0134190440)"
	if got != want {
		t.Errorf("Describe() = %q, want %q", got, want)
	}
}

var errBoom = errors.New("boom: repository unavailable")

func TestServiceDescribeError(t *testing.T) {
	fake := &fakeRepository{err: errBoom}
	svc := NewService(fake)

	_, err := svc.Describe("anything")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, errBoom) {
		t.Errorf("Describe() error = %v, want it to wrap %v", err, errBoom)
	}
}
