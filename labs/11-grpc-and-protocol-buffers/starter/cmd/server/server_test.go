// STARTER -- ILLUSTRATIVE, NOT GENERATOR OUTPUT.
//
// This test calls GetBooking directly as a plain Go function call --
// no grpc.Dial, no listener, no real network transport. Every
// generated gRPC method takes ctx context.Context first and is
// otherwise an ordinary Go method, so it's testable exactly like any
// other method. See Exercise 7 in exercise.md.
package main

import "testing"

// When you implement these, you'll need "context" (for context.Background())
// and "example.com/grpc-protobuf/bookingpb" (for bookingpb.GetBookingRequest)
// -- add them to the import block above once you use them.

func TestGetBooking_Found(t *testing.T) {
	t.Skip("TODO: implement TestGetBooking_Found")
}

func TestGetBooking_NotFound(t *testing.T) {
	t.Skip("TODO: implement TestGetBooking_NotFound")
}
