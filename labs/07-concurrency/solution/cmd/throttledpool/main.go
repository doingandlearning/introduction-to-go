// Exercise 9 (optional/stretch): the same worker pool from Exercise 8,
// extended with two more concurrency primitives: a lock-free counter from
// sync/atomic, and a time.Ticker that throttles how fast jobs reach the
// workers.
//
// Task 4 answer: time.Tick(...) returns a channel backed by a *Ticker
// that nothing ever holds a reference to, so nothing can ever call
// Stop() on it. In a program that exits on its own (like this one),
// that's harmless — the whole process tears down and takes the ticker
// with it. Inside a long-running service, calling a function that uses
// time.Tick repeatedly (once per request, once per loop iteration, etc.)
// leaks one running ticker per call, forever, because none of them can
// ever be stopped. time.NewTicker paired with a deferred Stop() is the
// version that's safe to use anywhere except a short, one-shot program.
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	numWorkers = 3
	numJobs    = 12
	ratePerSec = 5
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup, delivered *atomic.Int64) {
	defer wg.Done()
	for j := range jobs {
		fee := j * 2
		delivered.Add(1) // lock-free — no mutex needed for a single counter
		fmt.Printf("worker %d completed job %d, fee %d (total delivered: %d)\n",
			id, j, fee, delivered.Load())
		results <- fee
	}
}

func main() {
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup
	var delivered atomic.Int64

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg, &delivered)
	}

	start := time.Now()

	// Throttle: a real dispatch radio can only take ratePerSec updates a
	// second before the base station starts dropping packets. A
	// *time.Ticker, stopped explicitly, gates how fast jobs go out —
	// regardless of how fast the workers could otherwise drain them.
	ticker := time.NewTicker(time.Second / ratePerSec)
	defer ticker.Stop() // never leave a ticker running past the point you need it

	for j := 1; j <= numJobs; j++ {
		<-ticker.C // wait for the next tick before sending the next job
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	for r := range results {
		fmt.Println("completed delivery, fee:", r)
	}

	elapsed := time.Since(start)
	fmt.Printf("processed %d jobs at %d/sec in %s (total delivered: %d)\n",
		numJobs, ratePerSec, elapsed.Round(time.Millisecond), delivered.Load())
}
