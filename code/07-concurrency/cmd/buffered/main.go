// Command buffered is the same shape as cmd/producer, but the channel has
// capacity 3. It prints the send count immediately before each send so you
// can see exactly which one is the first to block.
//
//	go run ./cmd/buffered
//
// Compare against cmd/producer: with capacity 3, the first three sends
// return immediately (there's room in the buffer). The fourth send blocks
// until main receives at least one value, because the buffer is full.
package main

import "fmt"

func produce(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		fmt.Println("about to send", i)
		ch <- i // sends 1-3 don't block, send 4 blocks until room frees up
		fmt.Println("sent", i)
	}
	close(ch)
}

func main() {
	ch := make(chan int, 3) // buffered, capacity 3

	go produce(ch)

	for v := range ch {
		fmt.Println("received:", v)
	}

	fmt.Println("channel closed, done")
}
