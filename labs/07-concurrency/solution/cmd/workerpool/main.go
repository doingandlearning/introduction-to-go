// Exercise 8 (main lab exercise): a fixed pool of 3 dispatcher goroutines
// pulling delivery jobs off a shared jobs channel and pushing completion
// records onto a shared results channel.
package main

import (
	"fmt"
	"sync"
)

const numWorkers = 3

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fee := j * 2
		fmt.Printf("worker %d completed job %d, fee %d\n", id, j, fee)
		results <- fee
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= 12; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()      // every worker has drained jobs and exited
	close(results) // now safe — nobody is still sending to it

	for r := range results {
		fmt.Println("completed delivery, fee:", r)
	}
}
