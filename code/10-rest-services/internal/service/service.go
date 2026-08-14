// Package service owns business rules: validation, and whatever domain
// logic sits between "a request arrived" and "the store changed." It has
// no idea HTTP exists — no http.Request, no http.ResponseWriter, no
// status codes anywhere in this file.
package service

import (
	"errors"
	"fmt"

	"example.com/rest-service/internal/domain"
	"example.com/rest-service/internal/repository"
)

// ErrValidation is returned when caller-supplied data fails a business
// rule. The handler layer maps this to 400 Bad Request; this layer just
// reports what went wrong.
var ErrValidation = errors.New("validation error")

// ErrNotFound is re-exported so callers of this package (namely the
// handler) don't need to import the repository package just to compare
// errors with errors.Is.
var ErrNotFound = repository.ErrNotFound

// ItemService owns validation and business rules for items.
type ItemService struct {
	repo repository.ItemRepository
}

// NewItemService wires a service to whatever ItemRepository it's handed.
// This is constructor injection: the service depends on the
// repository.ItemRepository interface, not on repository.InMemoryRepository
// directly, so main can swap in a Postgres-backed repository later
// without this file changing at all — the same dependency-injection
// pattern introduced in Topics 8-9, now load-bearing in a real service.
func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

// Create validates and stores a new item.
func (s *ItemService) Create(item domain.Item) (domain.Item, error) {
	if err := validate(item); err != nil {
		return domain.Item{}, err
	}
	return s.repo.Create(item)
}

// Get fetches a single item by ID.
func (s *ItemService) Get(id string) (domain.Item, error) {
	return s.repo.Get(id)
}

// List returns every item.
func (s *ItemService) List() ([]domain.Item, error) {
	return s.repo.List()
}

// Update validates and replaces an existing item.
func (s *ItemService) Update(item domain.Item) (domain.Item, error) {
	if err := validate(item); err != nil {
		return domain.Item{}, err
	}
	return s.repo.Update(item)
}

// Delete removes an item by ID.
func (s *ItemService) Delete(id string) error {
	return s.repo.Delete(id)
}

// validate holds the actual business rule: a quantity that isn't
// positive isn't a real stock item. This is the kind of check that has
// nothing to do with HTTP and everything to do with the domain — exactly
// why it lives here and not in the handler.
func validate(item domain.Item) error {
	if item.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if item.Quantity <= 0 {
		return fmt.Errorf("%w: quantity must be positive", ErrValidation)
	}
	return nil
}
