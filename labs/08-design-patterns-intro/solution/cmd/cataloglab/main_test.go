package main_test

import "example.com/patterns-lab/internal/catalog"

// fakeRepository is a hand-written test double - no framework involved.
type fakeRepository struct{}

func (fakeRepository) FindByISBN(isbn string) (*catalog.Book, error) {
	return &catalog.Book{ISBN: isbn, Title: "Fake Book (test double)"}, nil
}
