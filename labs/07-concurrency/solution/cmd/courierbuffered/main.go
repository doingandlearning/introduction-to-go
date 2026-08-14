// Exercise 3: same scenario as Exercise 2, but over a BUFFERED channel
// (capacity 3), so the dispatcher can get ahead of the courier.
package main

import "fmt"

func dispatchOrders(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		fmt.Println("about to send order", i)
		ch <- i
	}
	close(ch)
}

// Answer: with capacity 3, sends 1, 2, and 3 all succeed immediately
// because there's room in the buffer. Send 4 is the first one that
// blocks — it has to wait until the courier receives at least one value,
// freeing a slot. From then on, each further send waits for a
// corresponding receive, same as the unbuffered version.

func main() {
	ch := make(chan int, 3) // buffered, capacity 3

	go dispatchOrders(ch)

	for orderID := range ch {
		fmt.Println("delivering order", orderID)
	}
}
