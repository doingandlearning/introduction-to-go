// ============================================================
// ILLUSTRATIVE -- NOT GENERATOR OUTPUT.
// ============================================================
//
// This file shows the shape of a gRPC client once real .pb.go files
// exist: dial a connection, get a generated client, call methods on
// it that read exactly like ordinary synchronous function calls --
// even though each one is a real network round trip.
//
// Because this sandbox has no google.golang.org/grpc and no network,
// the "connection" here is a fake, in-process stand-in
// (inProcessConn) that calls straight through to a server value in
// memory. The call SHAPE below -- client.GetBooking(ctx, req),
// stream, Recv() until io.EOF -- is exactly what you'd write against
// the real thing. Only the plumbing underneath (grpc.Dial /
// grpc.ClientConn) is faked out.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"example.com/grpc-protobuf/bookingpb"
)

// inProcessConn is a stand-in for a real *grpc.ClientConn. A real
// client obtains one via:
//
//   conn, err := grpc.Dial("localhost:50051",
//       grpc.WithTransportCredentials(insecure.NewCredentials()))
//   client := bookingpb.NewBookingServiceClient(conn)
//
// Here, we just hold a direct reference to a server implementation so
// the call-site code below can be written and read exactly as it
// would be against the real generated client.
type inProcessConn struct {
	srv bookingpb.BookingServiceServer
}

func (c *inProcessConn) GetBooking(ctx context.Context, req *bookingpb.GetBookingRequest) (*bookingpb.Booking, error) {
	return c.srv.GetBooking(ctx, req)
}

func (c *inProcessConn) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (bookingpb.BookingService_ListBookingsClient, error) {
	ch := make(chan *bookingpb.Booking, 16)
	go func() {
		defer close(ch)
		_ = c.srv.ListBookings(req, &channelStream{ch: ch})
	}()
	return &channelClientStream{ch: ch}, nil
}

type channelStream struct{ ch chan<- *bookingpb.Booking }

func (s *channelStream) Send(b *bookingpb.Booking) error {
	s.ch <- b
	return nil
}

type channelClientStream struct{ ch <-chan *bookingpb.Booking }

func (s *channelClientStream) Recv() (*bookingpb.Booking, error) {
	b, ok := <-s.ch
	if !ok {
		return nil, io.EOF
	}
	return b, nil
}

func main() {
	ctx := context.Background()

	// Real client setup would be:
	//   conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//   client := bookingpb.NewBookingServiceClient(conn)
	var client bookingpb.BookingServiceClient = &inProcessConn{srv: demoServer()}

	// --- Unary call: reads like an ordinary function call ---
	resp, err := client.GetBooking(ctx, &bookingpb.GetBookingRequest{Id: "abc123"})
	if err != nil {
		fmt.Println("GetBooking error:", err)
	} else {
		fmt.Printf("GetBooking: %+v\n", resp)
	}

	// --- Server-streaming call: read until io.EOF ---
	stream, err := client.ListBookings(ctx, &bookingpb.ListBookingsRequest{})
	if err != nil {
		fmt.Println("ListBookings error:", err)
		return
	}
	for {
		b, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("stream error:", err)
			break
		}
		fmt.Printf("ListBookings item: %+v\n", b)
	}
}

// demoServer wires up a tiny in-memory server so this file is
// self-contained. In cmd/server/main.go the equivalent type is called
// `server` -- duplicated here only so this client file can run on its
// own without importing an unexported type from another package.
func demoServer() bookingpb.BookingServiceServer { return &demoImpl{} }

type demoImpl struct{}

func (d *demoImpl) GetBooking(ctx context.Context, req *bookingpb.GetBookingRequest) (*bookingpb.Booking, error) {
	if req.Id != "abc123" {
		return nil, fmt.Errorf("booking %s not found", req.Id)
	}
	return &bookingpb.Booking{Id: "abc123", Guest: "A. Ortiz", Nights: 3}, nil
}

func (d *demoImpl) ListBookings(req *bookingpb.ListBookingsRequest, stream bookingpb.BookingService_ListBookingsServer) error {
	items := []*bookingpb.Booking{
		{Id: "abc123", Guest: "A. Ortiz", Nights: 3},
		{Id: "def456", Guest: "S. Chen", Nights: 1},
	}
	for _, b := range items {
		if err := stream.Send(b); err != nil {
			return err
		}
	}
	return nil
}
