// Command producer demonstrates the classic producer/consumer pattern over
// an UNBUFFERED channel: one goroutine sends five numbers, main receives
// and prints them with range, and the producer closes the channel when
// it's done sending.
//
//	go run ./cmd/producer
//
// Because the channel is unbuffered, each send blocks until main is ready
// to receive it — producer and consumer stay in lockstep, one value at a
// time.
package main

import (
	"fmt"
	"time"
)

// produce sends 1..5 down ch, one at a time, then closes it. Closing is a
// "no more values are coming" signal — it is always the sender's job, never
// the receiver's.
func produce(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i // blocks here until main receives
	}

	close(ch)
}

func consume(ch <-chan int, number int) {
	for v := range ch { // blocks here until produce sends a value
		fmt.Printf("consumer %d got: %d\n", number, v)
	}
}

// Goroutines -> large pieces of work!
// Streaming -> chunk by chunk

func main() {
	ch := make(chan int, 10) // [x] [x] [x] <- buffered channel, can hold 10 values before blocking

	go produce(ch)

	// range over a channel receives values until the channel is closed,
	// then exits the loop automatically — no manual "am I done" check.
	go consume(ch, 1)
	go consume(ch, 2)
	go consume(ch, 3)
	go consume(ch, 4)
	go consume(ch, 5)

	time.Sleep(1 * time.Second) // wait for consumers to finish
	fmt.Println("channel closed, done")
}
