// Exercise 9 (optional/stretch): the same worker pool from Exercise 8,
// extended with two more concurrency primitives: a lock-free counter from
// sync/atomic, and a time.Ticker that throttles how fast jobs reach the
// workers. The worker pool wiring itself is already done below — the
// TODOs are only the two new parts.
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numWorkers = 3
	numJobs    = 12
	ratePerSec = 5
)

// TODO(Exercise 9, task 1): add an *atomic.Int64 parameter named
// delivered, and call delivered.Add(1) once per completed job — no mutex
// needed for a single counter like this one.
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fee := j * 2
		fmt.Printf("worker %d completed job %d, fee %d\n", id, j, fee)
		results <- fee
	}
}

// TODO(Exercise 9, task 4): time.Tick(...) offers the same ticking
// channel without a *Ticker you can Stop(). Why would reaching for it
// here instead of time.NewTicker be a real problem in a long-running
// service, even though it works fine in this short-lived program?
func main() {
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// TODO(Exercise 9, task 1): declare `var delivered atomic.Int64` and
	// pass &delivered into each worker below.

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	start := time.Now()

	// TODO(Exercise 9, task 2): create a *time.Ticker firing every
	// time.Second/ratePerSec, defer its Stop(), and receive from
	// ticker.C once immediately before each send below so jobs go out no
	// faster than ratePerSec per second.

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	for r := range results {
		fmt.Println("completed delivery, fee:", r)
	}

	elapsed := time.Since(start)
	// TODO(Exercise 9, task 3): print delivered.Load() alongside elapsed
	// once the counter exists.
	fmt.Printf("processed %d jobs in %s\n", numJobs, elapsed.Round(time.Millisecond))
}
