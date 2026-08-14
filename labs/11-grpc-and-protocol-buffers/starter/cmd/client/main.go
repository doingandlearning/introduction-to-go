// STARTER -- ILLUSTRATIVE, NOT GENERATOR OUTPUT.
//
// TODO (Exercise 3): once your server is running (against real
// generated code), dial it with grpc.Dial and call GetBooking through
// the generated client stub. Time the call to confirm to yourself
// it's a real network round trip, even though it reads like an
// ordinary function call.
package main

import "fmt"

func main() {
	// TODO:
	//   conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//   client := bookingpb.NewBookingServiceClient(conn)
	//   resp, err := client.GetBooking(ctx, &bookingpb.GetBookingRequest{Id: "abc123"})
	fmt.Println("TODO (Exercise 3): implement the client call above")
}
