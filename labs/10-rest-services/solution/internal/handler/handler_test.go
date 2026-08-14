package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/rest-service/internal/domain"
	"example.com/rest-service/internal/repository"
	"example.com/rest-service/internal/service"
)

// newTestHandler wires a real handler over a real service over a real
// in-memory repository — the whole stack, minus any actual network.
func newTestHandler() *ItemHandler {
	repo := repository.NewInMemoryRepository()
	svc := service.NewItemService(repo)
	return NewItemHandler(svc)
}

// TestItemHandler_Create_ReturnsCreated drives the handler directly with
// httptest.NewRequest and httptest.NewRecorder — no mux, no
// ListenAndServe, no open port, no real network call.
func TestItemHandler_Create_ReturnsCreated(t *testing.T) {
	h := newTestHandler()

	body, err := json.Marshal(domain.Item{Name: "widget", Quantity: 5})
	if err != nil {
		t.Fatalf("marshaling request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var got domain.Item
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decoding response body: %v", err)
	}
	if got.Name != "widget" {
		t.Errorf("Name = %q, want %q", got.Name, "widget")
	}
	if got.Quantity != 5 {
		t.Errorf("Quantity = %d, want 5", got.Quantity)
	}
	if got.ID == "" {
		t.Error("expected a generated ID in the response, got empty string")
	}
}

// TestItemHandler_Create_ReturnsBadRequestOnInvalidQuantity confirms the
// handler maps a service validation error to 400, without any HTTP
// server actually running.
func TestItemHandler_Create_ReturnsBadRequestOnInvalidQuantity(t *testing.T) {
	h := newTestHandler()

	body, err := json.Marshal(domain.Item{Name: "widget", Quantity: 0})
	if err != nil {
		t.Fatalf("marshaling request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
