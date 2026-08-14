// Package service owns business rules: validation, and whatever domain
// logic sits between "a request arrived" and "the store changed." No
// HTTP types belong anywhere in this file.
package service

import (
	"errors"
	"fmt"

	"example.com/rest-service/internal/domain"
	"example.com/rest-service/internal/repository"
)

// ErrValidation is returned when caller-supplied data fails a business
// rule.
var ErrValidation = errors.New("validation error")

// ErrNotFound is re-exported so callers of this package don't need to
// import the repository package just to compare errors.
var ErrNotFound = repository.ErrNotFound

// ItemService owns validation and business rules for items.
type ItemService struct {
	repo repository.ItemRepository
}

// NewItemService wires a service to whatever ItemRepository it's handed
// — constructor injection, the same DI pattern from Topics 8-9.
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

// validate holds the actual business rule for an item.
func validate(item domain.Item) error {
	if item.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if item.Quantity <= 0 {
		return fmt.Errorf("%w: quantity must be positive", ErrValidation)
	}
	return nil
}
