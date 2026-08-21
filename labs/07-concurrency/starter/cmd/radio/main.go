// Exercise 5: select with a timeout. Dispatch waits on a status update
// from either of two couriers, but gives up after 2 seconds.
package main

import "fmt"

func main() {
	courierA := make(chan string)
	courierB := make(chan string)
	// Neither channel is ever sent on in this program — the timeout case
	// is guaranteed to fire. That's the point: we're building the pattern,
	// not simulating a real courier response.

	fmt.Println("waiting for a status update from courier A or B...")

	// TODO(Exercise 5): write a select with:
	//   - a case receiving from courierA
	//   - a case receiving from courierB
	//   - a case on time.After(2 * time.Second) that prints
	//     "no response from either courier, escalating"
}
