// Package catalog is Exercise 4's Service + Repository pair, built with
// manual constructor injection.
package catalog

import "fmt"

// Book is the domain type moved between layers.
type Book struct {
	ISBN  string
	Title string
}

// Repository is the interface Service depends on.
//
// TODO (Exercise 4): declare one method here:
//
//	FindByISBN(isbn string) (*Book, error)
type Repository interface {
	// TODO: add the FindByISBN method signature.
}

// Service owns lookup logic for the catalog. It depends on the
// Repository interface, not a concrete type.
type Service struct {
	repo Repository
}

// NewService is the constructor / injection point.
//
// TODO (Exercise 4): implement this. It should build and return a
// *Service with repo stored in its repo field.
func NewService(repo Repository) *Service {
	// TODO: replace this placeholder.
	return nil
}

// Describe fetches a book and formats a one-line description.
//
// TODO (Exercise 4): implement this. It should:
//  1. Call s.repo.FindByISBN(isbn).
//  2. If it errors, wrap and return the error (fmt.Errorf with %w).
//  3. On success, return a string like "Title (ISBN: 1234567890)".
func (s *Service) Describe(isbn string) (string, error) {
	// TODO: replace this placeholder.
	return "", fmt.Errorf("not implemented")
}

// InMemoryRepository is a real Repository implementation backed by a
// map. Already implemented - use it as-is for the "production" wiring
// in Exercise 4.
type InMemoryRepository struct {
	books map[string]*Book
}

// NewInMemoryRepository builds a repository seeded with the given books.
func NewInMemoryRepository(books map[string]*Book) *InMemoryRepository {
	return &InMemoryRepository{books: books}
}

// FindByISBN implements Repository.
func (r *InMemoryRepository) FindByISBN(isbn string) (*Book, error) {
	b, ok := r.books[isbn]
	if !ok {
		return nil, fmt.Errorf("book %q not found", isbn)
	}
	return b, nil
}
