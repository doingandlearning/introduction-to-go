// Command selectdemo shows select waiting on two independent input
// channels plus a time.After timeout case. Two goroutines "fetch" from
// warehouse-A and warehouse-B after a random delay; select reports each
// result as it arrives (in whichever order they actually finish), and
// gives up if either result takes longer than the timeout.
//
//	go run ./cmd/selectdemo
//
// Expected output (order and exact delays vary each run):
//
//	got: warehouse-A responded after 245ms
//	got: warehouse-B responded after 763ms
//	both warehouses reported in
package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// fetch simulates calling a real, independent data source (an API, a
// database, whatever) that takes a random amount of time to respond.
func fetch(source string, ch chan<- string) {
	delay := time.Duration(rand.IntN(1500)) * time.Millisecond
	time.Sleep(delay)
	ch <- fmt.Sprintf("%s responded after %v", source, delay)
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go fetch("warehouse-A", ch1)
	go fetch("warehouse-B", ch2)

	exit := false
	for !exit {
		select {
		case msg := <-ch1:
			fmt.Println("got:", msg)
		case msg := <-ch2:
			fmt.Println("got:", msg)
		case <-time.After(2 * time.Second):
			fmt.Println("timed out waiting for the rest")
			exit = true
		}
	}
	fmt.Println("both warehouses reported in")
}
