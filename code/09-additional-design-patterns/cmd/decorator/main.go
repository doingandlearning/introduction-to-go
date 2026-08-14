// Command decorator demonstrates the Decorator pattern in its most common
// Go shape: HTTP middleware. This exact structure - wrapping an
// http.Handler with another http.Handler - is the load-bearing idiom
// behind Topic 10's REST service, not just a textbook example.
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func baseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from the base handler"))
	})
}

// LoggingMiddleware wraps a handler and logs before and after it runs.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request:", r.URL.Path)
		next.ServeHTTP(w, r)
		log.Println("done")
	})
}

// HeaderMiddleware wraps a handler and stamps a response header before
// handing control to whatever it wraps.
func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "decorator-demo")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Stack both decorators around the same base handler. Order matters:
	// LoggingMiddleware is outermost here, so it logs the request even if
	// something further inside the chain misbehaves. Swap the nesting and
	// the observable behavior changes.
	handler := LoggingMiddleware(HeaderMiddleware(baseHandler()))

	req := httptest.NewRequest(http.MethodGet, "/guests", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	fmt.Println("status:", rec.Code)
	fmt.Println("X-Powered-By header:", rec.Header().Get("X-Powered-By"))
	fmt.Println("body:", rec.Body.String())
}
