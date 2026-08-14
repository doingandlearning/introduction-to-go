// Exercise 1: launch 5 couriers concurrently and wait for all of them to
// check in before the dispatcher starts the day.
package main

import "fmt"

func checkIn(courierID int) {
	fmt.Println("courier", courierID, "checked in")
}

func main() {
	// TODO(Exercise 1):
	//   1. Launch 5 goroutines, each calling checkIn(id) for a distinct id.
	//   2. Use a sync.WaitGroup so main waits for all 5 before printing
	//      "all couriers checked in, starting the day".
	//
	// Hint: wg.Add(1) before each go statement, defer wg.Done() inside the
	// goroutine, wg.Wait() before the final print.

	fmt.Println("all couriers checked in, starting the day")
}
