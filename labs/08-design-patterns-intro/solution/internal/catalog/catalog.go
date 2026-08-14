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
type Repository interface {
	FindByISBN(isbn string) (*Book, error)
}

// Service owns lookup logic for the catalog. It depends on the
// Repository interface, not a concrete type.
type Service struct {
	repo Repository
}

// NewService is the constructor / injection point.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Describe fetches a book and formats a one-line description.
func (s *Service) Describe(isbn string) (string, error) {
	b, err := s.repo.FindByISBN(isbn)
	if err != nil {
		return "", fmt.Errorf("describe: %w", err)
	}
	return fmt.Sprintf("%s (ISBN: %s)", b.Title, b.ISBN), nil
}

// InMemoryRepository is a real Repository implementation backed by a
// map.
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
