// Demo: "The struct tag footgun" — run this live, don't just describe it.
//
// Step 1: run it as-is, POST a correctly-spelled "name" field, watch it
// echo back correctly.
//
// Step 2: change the json tag below from `json:"name"` to `json:"geust"`
// (a "typo"), save, re-run, POST the SAME request again. Watch Name come
// back empty — no error, no warning, nothing in the logs. The compiler
// never checked the tag was spelled right; encoding/json just silently
// never found a matching key.
//
// Step 3: fix the tag back and show it populate again.
//
// Run:  go run ./cmd/struct-tag-footgun
// Hit:  curl -X POST localhost:8080/echo -d '{"name":"Ada","quantity":3}'
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Item — flip the json tag on Name between "name" and "geust" to run the
// footgun demo. Nothing else in this file needs to change.
type Item struct {
	Name     string `json:"name"` // <- change to `json:"geust"` for the demo
	Quantity int    `json:"quantity"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /echo", func(w http.ResponseWriter, r *http.Request) {
		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		log.Printf("decoded: %+v", item) // watch Name come back "" in the terminal

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
