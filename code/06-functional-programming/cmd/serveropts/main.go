// Command serveropts demonstrates the functional options pattern: a
// constructor that takes a variadic list of option functions instead of
// relying on overloading or default parameters, neither of which Go has.
//
//	go run ./cmd/serveropts
package main

import (
	"fmt"
	"time"
)

// Server is the thing being constructed. Its fields are unexported, so
// the only way to configure it from outside this package is through the
// options below - callers never touch the struct fields directly.
type Server struct {
	timeout time.Duration
	retries int
}

// ServerOption mutates a *Server. Every WithXxx function below returns
// one of these.
type ServerOption func(*Server)

// WithTimeout overrides the default timeout.
func WithTimeout(d time.Duration) ServerOption {
	return func(s *Server) { s.timeout = d }
}

// WithRetries overrides the default retry count.
func WithRetries(n int) ServerOption {
	return func(s *Server) { s.retries = n }
}

// NewServer builds a Server with sensible defaults, then applies
// whichever options were passed in, in order.
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		timeout: 30 * time.Second,
		retries: 3,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func main() {
	defaults := NewServer()
	fmt.Printf("defaults:      %+v\n", *defaults)

	customTimeout := NewServer(WithTimeout(5 * time.Second))
	fmt.Printf("custom timeout: %+v\n", *customTimeout)

	stacked := NewServer(WithTimeout(5*time.Second), WithRetries(10))
	fmt.Printf("stacked:       %+v\n", *stacked)
}
