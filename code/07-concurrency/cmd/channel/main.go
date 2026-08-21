package main

import "sync"

func main() {
	ch := make(chan int)
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		ch <- 42 // blocks here until the other goroutine receives
		ch <- 43
		close(ch)
		wg.Done()
	}()

	go func() {
		for v := range ch {
			println(v)
		}
		wg.Done()
	}()

	wg.Wait()
}
