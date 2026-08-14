// Command selectdemo shows select waiting on two input channels plus a
// time.After timeout case. Neither input channel is ever sent on, so the
// timeout case is guaranteed to fire.
//
//	go run ./cmd/selectdemo
//
// Expected output, after roughly a 2 second pause:
//
//	waiting on ch1, ch2, or a 2s timeout...
//	timed out waiting for ch1/ch2
package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	fmt.Println("waiting on ch1, ch2, or a 2s timeout...")

	select {
	case v := <-ch1:
		fmt.Println("from ch1:", v)
	case v := <-ch2:
		fmt.Println("from ch2:", v)
	case <-time.After(2 * time.Second):
		fmt.Println("timed out waiting for ch1/ch2")
	}
}
