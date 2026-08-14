// Command server runs the items REST API.
package main

import (
	"log"
	"net/http"

	"example.com/rest-service/internal/handler"
	"example.com/rest-service/internal/repository"
	"example.com/rest-service/internal/service"
)

func main() {
	// TODO(exercise 1): wire the composition root —
	// repo := repository.NewInMemoryRepository()
	// svc := service.NewItemService(repo)
	// h := handler.NewItemHandler(svc)
	_ = repository.NewInMemoryRepository
	_ = service.NewItemService
	_ = handler.NewItemHandler

	mux := http.NewServeMux()

	// TODO(exercise 1): add GET /ping, returning {"status":"ok"} as JSON.

	// TODO(exercise 2): mux.HandleFunc("POST /items", h.Create)
	// TODO(exercise 3): mux.HandleFunc("GET /items", h.List)
	// TODO(exercise 3): mux.HandleFunc("GET /items/{id}", h.Get)
	// TODO(exercise 4): mux.HandleFunc("PUT /items/{id}", h.Update)
	// TODO(exercise 4): mux.HandleFunc("DELETE /items/{id}", h.Delete)

	// TODO(exercise 5): wrap mux in handler.Logging(mux) below.
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
