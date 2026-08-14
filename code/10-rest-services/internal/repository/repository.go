// Package repository owns storage. It knows how to persist and fetch
// items — nothing about HTTP, nothing about validation rules.
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
// Depending on an interface here — not a concrete struct — is what lets
// main wire in an in-memory store today, and a Postgres-backed one later,
// without touching the service or handler layers at all.
type ItemRepository interface {
	Create(item domain.Item) (domain.Item, error)
	Get(id string) (domain.Item, error)
	List() ([]domain.Item, error)
	Update(item domain.Item) (domain.Item, error)
	Delete(id string) error
}

// InMemoryRepository is a map-backed ItemRepository guarded by a
// sync.RWMutex — direct continuity from Topic 7. Reads (Get, List) take
// the read lock and can run concurrently with each other; writes
// (Create, Update, Delete) take the exclusive write lock. Every request
// this repository serves is handled on its own goroutine, so this
// discipline isn't optional.
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
func (r *InMemoryRepository) Get(id string) (domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return domain.Item{}, ErrNotFound
	}
	return item, nil
}

// List returns every stored item. Order is not guaranteed — map
// iteration order in Go is deliberately randomized.
func (r *InMemoryRepository) List() ([]domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.Item, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	return items, nil
}

// Update replaces an existing item. Returns ErrNotFound if the ID isn't
// already present — this repository never silently creates on update.
func (r *InMemoryRepository) Update(item domain.Item) (domain.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[item.ID]; !ok {
		return domain.Item{}, ErrNotFound
	}
	r.items[item.ID] = item
	return item, nil
}

// Delete removes an item by ID. Returns ErrNotFound if it isn't there —
// deleting a nonexistent ID is a 404, never a panic.
func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}
