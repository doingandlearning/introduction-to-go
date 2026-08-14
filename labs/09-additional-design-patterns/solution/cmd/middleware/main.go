// Command middleware is Exercise 3: chain two HTTP middlewares around a
// base handler and confirm they run in the order you stacked them.
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func baseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ticket confirmed"))
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("LOG request:", r.URL.Path)
		next.ServeHTTP(w, r)
		log.Println("LOG done")
	})
}

func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Ticket-System", "lab9")
		next.ServeHTTP(w, r)
	})
}

func run(handler http.Handler, label string) {
	req := httptest.NewRequest(http.MethodGet, "/tickets/42", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	fmt.Println(label)
	fmt.Println("  status:", rec.Code)
	fmt.Println("  X-Ticket-System:", rec.Header().Get("X-Ticket-System"))
	fmt.Println("  body:", rec.Body.String())
}

func main() {
	logThenHeader := LoggingMiddleware(HeaderMiddleware(baseHandler()))
	run(logThenHeader, "-- LoggingMiddleware(HeaderMiddleware(base)) --")

	headerThenLog := HeaderMiddleware(LoggingMiddleware(baseHandler()))
	run(headerThenLog, "-- HeaderMiddleware(LoggingMiddleware(base)) --")

	// Exercise 3d: for this pair, the observable output (status, header,
	// body) is identical either way, because neither middleware can fail
	// or short-circuit the chain. The difference shows up once one
	// middleware can reject a request before the next runs (e.g. an auth
	// check) - then order changes what the inner handler ever sees.
}
