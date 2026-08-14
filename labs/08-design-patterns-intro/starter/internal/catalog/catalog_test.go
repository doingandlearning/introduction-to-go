// Exercise 5: prove dependency injection actually buys testability, by
// constructing Service with a fake Repository defined right here in the
// test - no mocking framework, no reflection, no real database.
package catalog

import "testing"

// TODO (Exercise 5): define a fakeRepository type here that implements
// Repository, i.e. a type with the method:
//
//	FindByISBN(isbn string) (*Book, error)
//
// Give it fields so a test can control what it returns, e.g.:
//
//	type fakeRepository struct {
//		book *Book
//		err  error
//	}
//
// If f.err is set, FindByISBN should return (nil, f.err); otherwise it
// should return (f.book, nil).

func TestServiceDescribeSuccess(t *testing.T) {
	t.Skip("TODO: implement TestServiceDescribeSuccess")
}

func TestServiceDescribeError(t *testing.T) {
	t.Skip("TODO: implement TestServiceDescribeError")
}
