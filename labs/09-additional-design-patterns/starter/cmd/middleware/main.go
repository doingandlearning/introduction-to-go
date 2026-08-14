// Command middleware is Exercise 3: chain two HTTP middlewares around a
// base handler and confirm they run in the order you stacked them.
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func baseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ticket confirmed"))
	})
}

// TODO (Exercise 3a): implement LoggingMiddleware(next http.Handler)
// http.Handler. It should append "LOG " to a package-level or closure
// variable (or just fmt.Println) before calling next.ServeHTTP.

// TODO (Exercise 3b): implement a second middleware, HeaderMiddleware,
// that sets a response header - e.g. w.Header().Set("X-Ticket-System",
// "lab9") - before calling next.ServeHTTP.

func main() {
	// TODO (Exercise 3c): chain both middlewares around baseHandler(),
	// invoke it with httptest.NewRequest / httptest.NewRecorder, and
	// print the status code, the X-Ticket-System header, and the body.

	// TODO (Exercise 3d): swap the nesting order and rerun. Write a
	// one-sentence comment here describing what, if anything, changed.

	fmt.Println("implement the TODOs above")
}
