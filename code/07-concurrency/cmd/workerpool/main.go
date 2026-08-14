// Command workerpool runs a small, fixed pool of worker goroutines that
// pull jobs off a shared `jobs` channel and push results onto a shared
// `results` channel. A sync.WaitGroup tracks when every worker has
// finished, which is what makes it safe to close `results`.
//
//	go run ./cmd/workerpool
package main

import (
	"fmt"
	"sync"
)

const numWorkers = 3

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs { // exits automatically when jobs is closed and drained
		fmt.Printf("worker %d squaring %d\n", id, j)
		results <- j * j
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	var wg sync.WaitGroup

	// Start a fixed crew of workers before any jobs exist.
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Hand out the work.
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs) // no more jobs coming — workers' range loops will end

	// Wait for every worker to finish draining jobs before closing results.
	// Closing results too early, while a worker is still sending to it,
	// would panic.
	wg.Wait()
	close(results)

	for r := range results {
		fmt.Println("result:", r)
	}
}
