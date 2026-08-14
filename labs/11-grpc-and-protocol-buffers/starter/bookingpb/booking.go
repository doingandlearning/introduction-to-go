// Package bookingpb -- STARTER, ILLUSTRATIVE STAND-IN, NOT GENERATOR OUTPUT.
//
// ============================================================
// Replace this whole file with the real booking.pb.go /
// booking_grpc.pb.go produced by protoc in Exercise 1, before doing
// anything else in this lab.
// ============================================================
//
// This stand-in exists only so the shape of GetBooking is visible
// before you've run protoc. It deliberately omits ListBookings --
// that's Exercise 4, which requires regenerating for real.
package bookingpb

import "context"

type Booking struct {
	Id     string
	Guest  string
	Nights int32
}

type GetBookingRequest struct {
	Id string
}

// BookingServiceServer -- the interface real generated code produces.
// A real version also gives you UnimplementedBookingServiceServer to
// embed for forward compatibility; omitted here since this is hand-written.
type BookingServiceServer interface {
	GetBooking(ctx context.Context, req *GetBookingRequest) (*Booking, error)
}

// BookingServiceClient -- the interface real generated code produces
// for callers.
type BookingServiceClient interface {
	GetBooking(ctx context.Context, req *GetBookingRequest) (*Booking, error)
}
