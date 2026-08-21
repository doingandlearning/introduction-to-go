// Exercise 1: launch 5 couriers concurrently and wait for all of them to
// check in before the dispatcher starts the day.
package main

import (
	"fmt"
	"sync"
)

func checkIn(courierID int) {
	fmt.Println("courier", courierID, "checked in")
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			checkIn(id)
		}(i) // IIFE - immediately invoked function expression
	}

	wg.Wait()
	fmt.Println("all couriers checked in, starting the day")
}
