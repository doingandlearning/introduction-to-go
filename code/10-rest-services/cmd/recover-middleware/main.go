// Demo: "One panic shouldn't take down every other request."
//
// Callback to Topic 2's "recover is reserved for truly exceptional
// situations" — a panicking handler is exactly that situation, and this
// is the real, small piece of production middleware that contains it.
//
// Run:   go run ./cmd/recover-middleware
// Step 1 (contained):
//
//	curl -i localhost:8080/panic   -> 500, logged, process still alive
//	curl -i localhost:8080/ping    -> still works — proves the process
//	                                   survived the panic above
//
// Step 2 (optional, do this as a SEPARATE run, in a terminal you don't
// mind losing): comment out the `handler.Recover` wrap below — i.e.
// change `Recover(mux)` to just `mux` — restart, hit /panic again, and
// watch the whole process die. That's the contrast: without the
// middleware, one bad request kills every other in-flight connection
// too, not just the one that panicked.
package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	})

	mux.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
		var items map[string]string // nil map
		items["boom"] = "this line panics: assignment to entry in nil map"
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", Recover(mux)))
}

// Recover wraps an http.Handler so a panic in any handler it serves gets
// caught, logged, and turned into a 500 — instead of crashing the whole
// process and taking every other in-flight request down with it.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("recovered panic: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
