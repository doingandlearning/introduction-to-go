// Exercise 2: producer/consumer over an UNBUFFERED channel. A dispatcher
// goroutine sends 5 order IDs one at a time; main receives and "delivers"
// each one.
package main

import "fmt"

func dispatchOrders(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i // blocks until main receives
	}
	close(ch)
}

func main() {
	ch := make(chan int) // unbuffered

	go dispatchOrders(ch)

	for orderID := range ch {
		fmt.Println("delivering order", orderID)
	}
}
