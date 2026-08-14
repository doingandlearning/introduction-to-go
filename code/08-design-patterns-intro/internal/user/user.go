// Package user demonstrates the Handler -> Service -> Repository
// layering pattern with manual constructor injection. There is no HTTP
// server here on purpose - the point of the pattern is that Service
// doesn't know or care whether it's called from a handler, a CLI, or a
// test.
package user

import "fmt"

// User is the domain type moved between layers.
type User struct {
	ID   string
	Name string
}

// Repository is the interface Service depends on. It knows nothing about
// HTTP or business rules - just "give me a user by ID."
type Repository interface {
	FindByID(id string) (*User, error)
}

// Service owns business logic. It depends on the Repository interface,
// not a concrete database type - that's what makes it testable with a
// fake, and it's dependency injection in its simplest Go form.
type Service struct {
	repo Repository
}

// NewService is the constructor. Whoever calls this decides what
// Repository implementation Service gets - a real one in production, a
// fake one in a test. This function is the injection point.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetGreeting is a small piece of "business logic": it fetches a user
// and formats a greeting. A real service would have validation and
// orchestration here too - this stands in for that.
func (s *Service) GetGreeting(id string) (string, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return "", fmt.Errorf("get greeting: %w", err)
	}
	return fmt.Sprintf("Hello, %s (id=%s)", u.Name, u.ID), nil
}

// InMemoryRepository is a "real" (if minimal) Repository implementation
// backed by a map. In production this might instead be backed by
// Postgres, Redis, an external API - Service would not change either way.
type InMemoryRepository struct {
	users map[string]*User
}

// NewInMemoryRepository builds a repository seeded with the given users.
func NewInMemoryRepository(users map[string]*User) *InMemoryRepository {
	return &InMemoryRepository{users: users}
}

// FindByID implements Repository.
func (r *InMemoryRepository) FindByID(id string) (*User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %q not found", id)
	}
	return u, nil
}
