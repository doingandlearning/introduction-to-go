// Demo: "What a REST service is" — the opening slide of Topic 10.
//
// Type this live, from scratch, in front of the room. The point is that
// this is the WHOLE thing — no project generator, no framework import,
// no app.get(...) wrapper. Just the standard library.
//
// Run:   go run ./cmd/minimal-ping
// Hit:   curl localhost:8080/ping
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "okay"})
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
