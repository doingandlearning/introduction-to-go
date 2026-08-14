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
//
// TODO(exercise 2): decode the request body into a domain.Item (400 on
// bad JSON), call h.svc.Create, map a service error with
// writeServiceError, and on success write 201 Created with the item
// JSON-encoded.
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement Create")
}

// Get handles GET /items/{id}.
//
// TODO(exercise 3): read r.PathValue("id"), call h.svc.Get, map errors
// with writeServiceError, otherwise write 200 with the item JSON-encoded.
func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement Get")
}

// List handles GET /items.
//
// TODO(exercise 3): call h.svc.List and write 200 with the items
// JSON-encoded.
func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement List")
}

// Update handles PUT /items/{id}.
//
// TODO(exercise 4): same shape as Create, but set item.ID from
// r.PathValue("id") before calling h.svc.Update, and return 200 (not
// 201) on success.
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement Update")
}

// Delete handles DELETE /items/{id}.
//
// TODO(exercise 4): read r.PathValue("id"), call h.svc.Delete, map
// errors with writeServiceError, otherwise write 204 No Content.
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement Delete")
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

// jsonHeader is a small helper so every success path sets the same
// Content-Type header the same way.
func jsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

var _ = json.Marshal // keep encoding/json imported until Create/Get/etc. use it directly
