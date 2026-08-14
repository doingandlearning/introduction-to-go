// ============================================================
// ILLUSTRATIVE -- NOT GENERATOR OUTPUT.
// ============================================================
//
// This file shows roughly what a hand-written gRPC server
// implementation looks like once real .pb.go files exist. The
// interface it implements (bookingpb.BookingServiceServer) is itself
// a hand-written approximation -- see bookingpb/booking.go for the
// full explanation.
//
// What's REAL Go here: the sync.RWMutex-protected in-memory store is
// a direct, working continuation of the Topic 10 REST pattern, and
// the logic itself (lock, look up, unlock) is exactly what you'd
// write against the real generated interface.
//
// What's MISSING to make this a real gRPC server: registering this
// implementation with an actual *grpc.Server, listening on a TCP
// port, and calling grpc.Serve(). That needs google.golang.org/grpc,
// which isn't available in this sandbox. The commented-out block at
// the bottom shows exactly what that wiring looks like in a real
// environment.
package main

import (
	"context"
	"fmt"
	"sync"

	"example.com/grpc-protobuf/bookingpb"
)

// server implements bookingpb.BookingServiceServer.
//
// In real generated code you would also embed
// bookingpb.UnimplementedBookingServiceServer here -- that's what lets
// your server keep compiling if a new rpc is added to the .proto later
// and you haven't implemented it yet. Omitted here since this file's
// BookingServiceServer is hand-written, not generated.
type server struct {
	mu    sync.RWMutex
	store map[string]*bookingpb.Booking
}

func newServer() *server {
	return &server{
		store: map[string]*bookingpb.Booking{
			"abc123": {Id: "abc123", Guest: "A. Ortiz", Nights: 3},
			"def456": {Id: "def456", Guest: "S. Chen", Nights: 1},
			"ghi789": {Id: "ghi789", Guest: "M. Osei", Nights: 5},
		},
	}
}

// GetBooking -- unary RPC. Notice the first parameter: context.Context.
// Every generated method takes one, first -- not a gRPC quirk, a
// Go-wide idiom for cancellation, deadlines, and request-scoped values.
func (s *server) GetBooking(ctx context.Context, req *bookingpb.GetBookingRequest) (*bookingpb.Booking, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.store[req.Id]
	if !ok {
		// A real implementation returns a gRPC status error here, e.g.
		// status.Errorf(codes.NotFound, "booking %s not found", req.Id),
		// using google.golang.org/grpc/codes and .../status. Simplified
		// to a plain error since this file can't import grpc-go.
		return nil, fmt.Errorf("booking %s not found", req.Id)
	}
	return b, nil
}

// ListBookings -- server-streaming RPC. Instead of returning one value,
// it loops over the store and calls stream.Send for each item, exactly
// as it would against the real generated stream type.
func (s *server) ListBookings(req *bookingpb.ListBookingsRequest, stream bookingpb.BookingService_ListBookingsServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, b := range s.store {
		if err := stream.Send(b); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	fmt.Println("This illustrative server has no real network transport --")
	fmt.Println("see the commented block below for what real wiring looks like.")
	_ = newServer()

	// ---- Real wiring, once .pb.go files exist (do not uncomment here) ----
	//
	// lis, err := net.Listen("tcp", ":50051")
	// if err != nil {
	//     log.Fatalf("failed to listen: %v", err)
	// }
	// grpcServer := grpc.NewServer()
	// bookingpb.RegisterBookingServiceServer(grpcServer, newServer())
	// log.Println("BookingService listening on :50051")
	// if err := grpcServer.Serve(lis); err != nil {
	//     log.Fatalf("failed to serve: %v", err)
	// }
}
