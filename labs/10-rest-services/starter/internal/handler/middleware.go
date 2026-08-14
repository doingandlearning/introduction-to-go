package handler

import (
	"net/http"
)

// Logging wraps an http.Handler and adds request logging — the
// Decorator pattern from Topic 9.
//
// TODO(exercise 5): return an http.HandlerFunc that records time.Now(),
// calls next.ServeHTTP(w, r), then logs the method, path, and elapsed
// duration with log.Printf.
func Logging(next http.Handler) http.Handler {
	panic("TODO: implement Logging")
}
