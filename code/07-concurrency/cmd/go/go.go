package main

import (
	"fmt"
	"sync"
)

func sayHi() {
	fmt.Println("hi")
}

func main() {
	go sayHi()          // runs concurrently
	fmt.Println("main") // main doesn't wait for sayHi
	// wait for 1 second
	batch()
}

// Goroutines (unit of work) - 2KB in Go vs 1MB in Java
// Machines (threads on the OS)
// Processors

// MAXPROC - 1

func batch() {
	var wg sync.WaitGroup
	batchSize := 10
	population := 1000
	for i := 0; i < population; i += batchSize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + batchSize
			if end > population {
				end = population
			}
			fmt.Printf("Processing %d to %d\n", start, end)
		}(i)
	}
	wg.Wait()
	fmt.Println("All done")
}
