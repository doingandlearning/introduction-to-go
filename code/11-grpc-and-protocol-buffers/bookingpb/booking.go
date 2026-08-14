// Package bookingpb is a HAND-WRITTEN APPROXIMATION of what
// `protoc --go_out=. --go-grpc_out=. booking.proto` would generate.
//
// ============================================================
// THIS IS NOT GENERATOR OUTPUT. DO NOT TREAT IT AS AUTHORITATIVE.
// ============================================================
//
// This file exists only so the illustrative server and client code in
// this directory has something to compile against conceptually, and so
// delegates can see roughly the SHAPE of what protoc produces before
// they run it for real. It is deliberately simplified:
//
//   - Real generated code embeds proto.Message plumbing (Reset, String,
//     ProtoReflect, wire marshal/unmarshal, etc.) on every message type.
//     None of that is here.
//   - Real generated code produces a raw gRPC service descriptor,
//     *_grpc.pb.go client/server interfaces wired to google.golang.org/grpc,
//     and an UnimplementedBookingServiceServer for forward compatibility.
//     This file fakes just enough of that shape to read naturally.
//   - This file will NOT actually run against a real gRPC transport --
//     there is no protoc, no grpc-go, and no network in this sandbox.
//
// BEFORE using any of this in a classroom or real project: delete this
// file and regenerate the real thing with protoc against booking.proto.
// See code/README.md for the exact commands.
package bookingpb

import "context"

// Booking -- in real generated code this is produced from the `Booking`
// message in booking.proto, and additionally implements proto.Message.
type Booking struct {
	Id     string
	Guest  string
	Nights int32
}

// GetBookingRequest -- generated from the message of the same name.
type GetBookingRequest struct {
	Id string
}

// ListBookingsRequest -- generated from the message of the same name.
// Empty today; protoc would still generate a concrete (empty) struct.
type ListBookingsRequest struct{}

// BookingServiceServer is the interface real generated code produces
// for you to implement -- one method per rpc in the .proto service.
//
// The real generated version also gives you
// UnimplementedBookingServiceServer, embedded in your implementation,
// so adding a new rpc to the .proto later doesn't break every existing
// server implementation immediately (forward compatibility for the
// interface itself).
type BookingServiceServer interface {
	GetBooking(ctx context.Context, req *GetBookingRequest) (*Booking, error)
	// ListBookings is server-streaming: instead of returning a value,
	// real generated code gives your implementation a stream object to
	// call .Send(...) on, once per item, for as long as the call is open.
	ListBookings(req *ListBookingsRequest, stream BookingService_ListBookingsServer) error
}

// BookingService_ListBookingsServer is the illustrative stand-in for the
// generated server-side stream handle. Real generated code wires this to
// grpc.ServerStream under the hood.
type BookingService_ListBookingsServer interface {
	Send(*Booking) error
}

// BookingServiceClient is the interface real generated code produces for
// callers -- this is what makes a network call "read like an ordinary
// function call" from the caller's point of view.
type BookingServiceClient interface {
	GetBooking(ctx context.Context, req *GetBookingRequest) (*Booking, error)
	ListBookings(ctx context.Context, req *ListBookingsRequest) (BookingService_ListBookingsClient, error)
}

// BookingService_ListBookingsClient is the illustrative stand-in for the
// generated client-side stream handle. Real generated code returns
// io.EOF from Recv() when the stream is done, exactly as shown in the
// hand-written client in cmd/client/main.go.
type BookingService_ListBookingsClient interface {
	Recv() (*Booking, error)
}
