// Command server runs the items REST API.
//
// main is the composition root: the one place in the program that knows
// about every concrete type. Everything it wires together only ever sees
// interfaces from that point on — repository.ItemRepository,
// service interface, handler.ItemService.
package main

import (
	"log"
	"net/http"

	"example.com/rest-service/internal/handler"
	"example.com/rest-service/internal/repository"
	"example.com/rest-service/internal/service"
)

func main() {
	repo := repository.NewInMemoryRepository()
	svc := service.NewItemService(repo)
	h := handler.NewItemHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("POST /items", h.Create)
	mux.HandleFunc("GET /items", h.List)
	mux.HandleFunc("GET /items/{id}", h.Get)
	mux.HandleFunc("PUT /items/{id}", h.Update)
	mux.HandleFunc("DELETE /items/{id}", h.Delete)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler.Logging(mux)))
}
