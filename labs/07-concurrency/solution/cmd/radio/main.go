// Exercise 5: select with a timeout. Dispatch waits on a status update
// from either of two couriers, but gives up after 2 seconds.
package main

import (
	"fmt"
	"time"
)

func main() {
	courierA := make(chan string)
	courierB := make(chan string)

	fmt.Println("waiting for a status update from courier A or B...")

	select {
	case v := <-courierA:
		fmt.Println("courier A:", v)
	case v := <-courierB:
		fmt.Println("courier B:", v)
	case <-time.After(2 * time.Second):
		fmt.Println("no response from either courier, escalating")
	}
}
