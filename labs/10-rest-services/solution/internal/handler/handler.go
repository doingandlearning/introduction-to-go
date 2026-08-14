// Package handler owns HTTP concerns only: decode a request, call the
// service, translate the result into a status code and a JSON body.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"example.com/rest-service/internal/domain"
	"example.com/rest-service/internal/service"
)

// ItemService is the contract the handler depends on.
type ItemService interface {
	Create(item domain.Item) (domain.Item, error)
	Get(id string) (domain.Item, error)
	List() ([]domain.Item, error)
	Update(item domain.Item) (domain.Item, error)
	Delete(id string) error
}

// ItemHandler translates HTTP requests into service calls and service
// results back into HTTP responses.
type ItemHandler struct {
	svc ItemService
}

// NewItemHandler wires a handler to whatever ItemService it's handed —
// constructor injection, same shape as service.NewItemService.
func NewItemHandler(svc ItemService) *ItemHandler {
	return &ItemHandler{svc: svc}
}

// Create handles POST /items.
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var item domain.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	created, err := h.svc.Create(item)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// Get handles GET /items/{id}.
func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	item, err := h.svc.Get(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// List handles GET /items.
func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List()
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Update handles PUT /items/{id}.
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var item domain.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	item.ID = id

	updated, err := h.svc.Update(item)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// Delete handles DELETE /items/{id}.
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.svc.Delete(id); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeServiceError maps every service error this handler can receive to
// a status code, explicitly, in one place.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, service.ErrValidation):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
