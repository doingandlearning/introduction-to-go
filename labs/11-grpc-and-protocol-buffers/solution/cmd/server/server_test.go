// ============================================================
// ILLUSTRATIVE -- NOT GENERATOR OUTPUT.
// ============================================================
//
// This test calls GetBooking directly as a plain Go function call --
// no grpc.Dial, no listener, no real network transport. That's not a
// limitation of this illustrative code: every generated gRPC method
// takes ctx context.Context first and is otherwise an ordinary Go
// method, so testing it is identical to testing any other method.
// Once Exercise 1's real generated types exist, this exact test still
// compiles and passes unchanged, because the method signature doesn't
// change.
package main

import (
	"context"
	"testing"

	"example.com/grpc-protobuf/bookingpb"
)

func TestGetBooking_Found(t *testing.T) {
	s := newServer()

	got, err := s.GetBooking(context.Background(), &bookingpb.GetBookingRequest{Id: "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Guest != "A. Ortiz" {
		t.Errorf("GetBooking(abc123).Guest = %q, want %q", got.Guest, "A. Ortiz")
	}
	if got.Nights != 3 {
		t.Errorf("GetBooking(abc123).Nights = %d, want 3", got.Nights)
	}
}

func TestGetBooking_NotFound(t *testing.T) {
	s := newServer()

	_, err := s.GetBooking(context.Background(), &bookingpb.GetBookingRequest{Id: "does-not-exist"})
	if err == nil {
		t.Fatal("expected an error for an unknown booking ID, got nil")
	}
}
