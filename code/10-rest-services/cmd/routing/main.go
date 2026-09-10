// Demo: "Go 1.22 changed the router."
//
// Three routes, same-looking path, no third-party router (gorilla/mux,
// chi) needed. Before Go 1.22 none of this method/path-param matching
// existed in net/http — you'd have written your own string-splitting or
// pulled in a dependency just for this.
//
// Run:   go run ./cmd/routing
// Hit:
//
//	curl localhost:8080/items/42            -> GET, prints the id
//	curl -X POST localhost:8080/items       -> different handler, same path
//	curl -X DELETE localhost:8080/items/42  -> different handler again
//	curl localhost:8080/items/42/history    -> 404, path doesn't match anything registered
//
// Worth saying out loud: check `go version` before trusting an older
// tutorial's routing setup — anything written before Feb 2024 assumes
// none of this exists yet.
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "GET item %s\n", id)
	})

	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "POST: create a new item")
	})

	mux.HandleFunc("DELETE /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "DELETE item %s\n", id)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
