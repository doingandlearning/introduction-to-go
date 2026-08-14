// Command waitgroup launches five goroutines and uses sync.WaitGroup to
// wait for all of them to finish before main exits.
//
// Run it several times in a row:
//
//	go run ./cmd/waitgroup
//
// The five "worker" lines print in a different order almost every run.
// WaitGroup guarantees that main waits for all five to complete — it says
// nothing at all about the order they finish in.
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Println("worker", n, "done")
		}(i)
	}

	wg.Wait() // blocks until all 5 goroutines have called wg.Done()
	fmt.Println("all workers finished")
}
