package handler

import (
	"log"
	"net/http"
	"time"
)

// Logging wraps an http.Handler and logs the method, path, and duration
// of every request that passes through it. This is the Decorator pattern
// from Topic 9 — a function taking an http.Handler and returning a new
// one that wraps it — and it's the exact shape used by essentially every
// piece of Go HTTP middleware you'll meet in the wild.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
