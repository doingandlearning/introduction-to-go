// Package repository owns storage.
package repository

import (
	"errors"
	"strconv"
	"sync"

	"example.com/rest-service/internal/domain"
)

// ErrNotFound is returned when a lookup, update, or delete targets an ID
// that doesn't exist in the store.
var ErrNotFound = errors.New("item not found")

// ItemRepository is the storage contract the service layer depends on.
type ItemRepository interface {
	Create(item domain.Item) (domain.Item, error)
	Get(id string) (domain.Item, error)
	List() ([]domain.Item, error)
	Update(item domain.Item) (domain.Item, error)
	Delete(id string) error
}

// InMemoryRepository is a map-backed ItemRepository guarded by a
// sync.RWMutex — direct continuity from Topic 7.
type InMemoryRepository struct {
	mu     sync.RWMutex
	items  map[string]domain.Item
	nextID int
}

// NewInMemoryRepository builds an empty, ready-to-use repository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: make(map[string]domain.Item),
	}
}

// Create assigns a new ID and stores the item.
func (r *InMemoryRepository) Create(item domain.Item) (domain.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	item.ID = strconv.Itoa(r.nextID)
	r.items[item.ID] = item
	return item, nil
}

// Get fetches a single item by ID.
//
// TODO(exercise 3): implement this. Use r.mu.RLock()/RUnlock() — it's a
// read. Return ErrNotFound if the ID isn't in the map.
func (r *InMemoryRepository) Get(id string) (domain.Item, error) {
	panic("TODO: implement Get")
}

// List returns every stored item.
//
// TODO(exercise 3): implement this. Use r.mu.RLock()/RUnlock() — it's a
// read.
func (r *InMemoryRepository) List() ([]domain.Item, error) {
	panic("TODO: implement List")
}

// Update replaces an existing item.
//
// TODO(exercise 4): implement this. Use r.mu.Lock()/Unlock() — it's a
// write. Return ErrNotFound if the ID isn't already present; never
// silently create on update.
func (r *InMemoryRepository) Update(item domain.Item) (domain.Item, error) {
	panic("TODO: implement Update")
}

// Delete removes an item by ID.
//
// TODO(exercise 4): implement this. Use r.mu.Lock()/Unlock() — it's a
// write. Return ErrNotFound if the ID isn't present; never panic on a
// missing ID.
func (r *InMemoryRepository) Delete(id string) error {
	panic("TODO: implement Delete")
}
