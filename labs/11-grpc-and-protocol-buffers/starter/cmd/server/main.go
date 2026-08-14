// STARTER -- ILLUSTRATIVE, NOT GENERATOR OUTPUT.
//
// Once you have real generated types (Exercise 1), implement the
// TODOs below against them. This file, as-is, shows the shape but
// does not implement the logic.
package main

import (
	"context"
	"fmt"
	"sync"

	"example.com/grpc-protobuf/bookingpb"
)

// server implements bookingpb.BookingServiceServer.
type server struct {
	mu    sync.RWMutex
	store map[string]*bookingpb.Booking
}

func newServer() *server {
	return &server{
		store: map[string]*bookingpb.Booking{
			"abc123": {Id: "abc123", Guest: "A. Ortiz", Nights: 3},
		},
	}
}

// TODO (Exercise 2): implement GetBooking.
//   - RLock/RUnlock around the map read (Topic 10's pattern)
//   - if the ID isn't found, return an error (a real version returns
//     status.Errorf(codes.NotFound, ...) from google.golang.org/grpc)
func (s *server) GetBooking(ctx context.Context, req *bookingpb.GetBookingRequest) (*bookingpb.Booking, error) {
	panic("TODO: implement GetBooking")
}

func main() {
	fmt.Println("TODO (Exercise 2): register this server with a real *grpc.Server")
	fmt.Println("and listen on :50051. See code/README.md for the wiring shape.")
	_ = newServer()
}
