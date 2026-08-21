package main

import (
	"fmt"
	"sync"
)

func increment(wg *sync.WaitGroup, incs chan<- int) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		incs <- 1 // send instead of touching counter
	}
}

func incrementer(incs <-chan int, result chan<- int) {
	total := 0
	for range incs {
		total++
	}
	result <- total
}

func main() {
	incs := make(chan int)
	result := make(chan int)

	go incrementer(incs, result) // single goroutine owns the counter

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go increment(&wg, incs)
	}

	wg.Wait()
	close(incs) // close the channel to signal the goroutine to finish
	fmt.Println("final counter value:", <-result)
}
