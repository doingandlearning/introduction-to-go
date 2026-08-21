// Exercise 8 (main lab exercise): a fixed pool of 3 dispatcher goroutines
// pulling delivery jobs off a shared jobs channel and pushing completion
// records onto a shared results channel.
package main

import (
	"fmt"
	"sync"
)

const numWorkers = 3

// TODO(Exercise 8): implement worker so that it:
//  1. ranges over jobs (this exits automatically once jobs is closed and
//     drained)
//  2. for each job j received, sends j*2 (the "delivery fee") to results
//  3. calls wg.Done() via defer when the function returns
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	// TODO
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	var wg sync.WaitGroup

	// TODO(Exercise 8):
	//   1. Launch exactly numWorkers workers, all sharing jobs and results.
	//   2. Send job IDs 1 through 12 into jobs, then close jobs.
	//   3. wg.Wait() for every worker to finish BEFORE closing results —
	//      closing results while a worker is still sending to it panics.
	//   4. Range over results and print each one.

	_ = jobs
	_ = results
	_ = &wg
	fmt.Println("TODO: implement the dispatch pool")
}
