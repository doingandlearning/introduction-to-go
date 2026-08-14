// Exercise 6: INTENTIONALLY BROKEN — run to observe, then fix it in the
// sibling files main_mutex.go and main_channel.go (task 3 and 4).
//
// Two dispatcher goroutines both increment a shared deliveryCount a
// thousand times each, with no synchronization. deliveryCount++ is a
// read-modify-write, not atomic, so increments get lost when the two
// goroutines interleave.
//
//	go run ./standalone/race          # "succeeds" with a wrong number
//	go run -race ./standalone/race    # reports the exact conflicting accesses
package main

import (
	"fmt"
	"sync"
)

var deliveryCount int

func recordDeliveries(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		deliveryCount++ // no synchronization — the race
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go recordDeliveries(&wg)
	}

	wg.Wait()
	fmt.Println("final delivery count:", deliveryCount)
}
