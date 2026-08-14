package handler

import (
	"log"
	"net/http"
	"time"
)

// Logging wraps an http.Handler and logs the method, path, and duration
// of every request that passes through it — the Decorator pattern from
// Topic 9.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
