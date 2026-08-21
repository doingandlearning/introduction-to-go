// Exercise 4: same scenario as Exercise 3, but over a BUFFERED channel
// (capacity 3), so the dispatcher can get ahead of the courier.
package main

// TODO(Exercise 4): add "fmt" to the import block above, then implement
// dispatchOrders so it sends order IDs 1 through 5 into ch, printing
// "about to send order <n>" immediately before each send, then closes ch
// when done.
func dispatchOrders(ch chan<- int) {
	// TODO
}

// TODO(Exercise 4, as a comment): which send is the first one that would
// block if the courier were slow to start receiving, and why? Write your
// answer here before moving on.

func main() {
	ch := make(chan int, 3) // buffered, capacity 3

	go dispatchOrders(ch)

	// TODO(Exercise 4): receive from ch with range and print
	// "delivering order <n>" for each value received.
}
