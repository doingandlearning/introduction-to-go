package service

import (
	"errors"
	"testing"

	"example.com/rest-service/internal/domain"
	"example.com/rest-service/internal/repository"
)

// TestItemService_Create_StoresValidItem exercises the service layer
// directly, with a real in-memory repository behind it and no HTTP
// involved anywhere.
func TestItemService_Create_StoresValidItem(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := NewItemService(repo)

	item := domain.Item{Name: "widget", Quantity: 5}

	created, err := svc.Create(item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID == "" {
		t.Error("expected a generated ID, got empty string")
	}
	if created.Name != "widget" {
		t.Errorf("Name = %q, want %q", created.Name, "widget")
	}
	if created.Quantity != 5 {
		t.Errorf("Quantity = %d, want 5", created.Quantity)
	}

	// Confirm it actually landed in the repository, not just returned.
	fetched, err := svc.Get(created.ID)
	if err != nil {
		t.Fatalf("Get after Create: unexpected error: %v", err)
	}
	if fetched != created {
		t.Errorf("Get(%q) = %+v, want %+v", created.ID, fetched, created)
	}
}

// TestItemService_Create_RejectsInvalidQuantity proves the validation
// rule without spinning up any HTTP machinery.
func TestItemService_Create_RejectsInvalidQuantity(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := NewItemService(repo)

	_, err := svc.Create(domain.Item{Name: "widget", Quantity: 0})
	if err == nil {
		t.Fatal("expected an error for quantity 0, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected errors.Is(err, ErrValidation), got %v", err)
	}
}
