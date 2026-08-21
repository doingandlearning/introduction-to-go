// Exercise 3: producer/consumer over an UNBUFFERED channel. A dispatcher
// goroutine sends 5 order IDs one at a time; main receives and "delivers"
// each one.
package main

// TODO(Exercise 3): implement dispatchOrders so it sends order IDs 1
// through 5 into ch, one at a time, then closes ch when done. Closing is
// the sender's responsibility.
func dispatchOrders(ch chan<- int) {
	// TODO
}

func main() {
	ch := make(chan int) // unbuffered

	go dispatchOrders(ch)

	// TODO(Exercise 3): add "fmt" to the import block above, then receive
	// from ch with range and print "delivering order <n>" for each value
	// received. range should exit on its own once dispatchOrders closes
	// the channel.
}
