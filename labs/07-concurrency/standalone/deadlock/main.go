// Exercise 6: INTENTIONALLY BROKEN — run to observe, this is the point of
// the exercise, not a bug to silently fix.
//
// Dispatch tries to hand a single urgent order directly to a courier over
// an unbuffered channel, but no courier goroutine was ever started to
// receive it. main blocks forever on the send.
//
//	go run ./standalone/deadlock
//
// Expected output:
//
//	fatal error: all goroutines are asleep - deadlock!
//
// Fill in your answer to Exercise 6, task 3, at the bottom of this file.
package main

import "fmt"

func main() {
	urgentOrder := make(chan int) // unbuffered

	fmt.Println("handing off urgent order...")
	urgentOrder <- 42 // nobody is receiving — this blocks forever

	fmt.Println("this line never runs")
}

// Exercise 6, task 3 — your answer:
//
// (Write one sentence here describing a scenario where a deadlock happens
// between only SOME of your goroutines, and explain why the runtime's
// "all goroutines are asleep" detector would not report it.)
