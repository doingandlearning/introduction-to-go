package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func fetch(source string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	delay := time.Duration(rand.IntN(1500)) * time.Millisecond
	time.Sleep(delay)
	ch <- fmt.Sprintf("%s responded after %v", source, delay)
}

func main() {
	results := make(chan string)
	var wg sync.WaitGroup

	sources := []string{"warehouse-A", "warehouse-B"}
	wg.Add(len(sources))
	for _, s := range sources {
		go fetch(s, results, &wg)
	}

	go func() {
		wg.Wait()      // wait until every fetch has sent its one result
		close(results) // only now is it guaranteed safe
	}()

	for msg := range results { // exits automatically once results is closed
		fmt.Println("got:", msg)
	}
	fmt.Println("both warehouses reported in")
}
