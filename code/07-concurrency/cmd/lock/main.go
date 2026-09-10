package main

import (
	"sync"
)

var mu sync.Mutex

type Counter int

func (c *Counter) Increment() {
	mu.Lock()
	defer mu.Unlock()
	*c++
}

var counter Counter

func increment() {
	for i := 0; i < 1000; i++ {
		counter.Increment()
	}
}

func wasteSomeTime() {
	for i := 0; i < 1000000; i++ {
		_ = i * i
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			increment()
		}()
	}

	go wasteSomeTime()

	wg.Wait()
	println("final counter value:", counter)
}
