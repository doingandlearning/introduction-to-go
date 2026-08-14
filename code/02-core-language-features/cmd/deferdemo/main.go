// Command deferdemo shows defer's LIFO ordering, and the classic gotcha
// where a loop variable is captured by reference to the deferred call's
// evaluation time rather than the value you expected.
//
//	go run ./cmd/deferdemo
package main

import "fmt"

func orderedDefers() {
	fmt.Println("--- orderedDefers ---")
	defer fmt.Println("deferred: first")
	defer fmt.Println("deferred: second")
	defer fmt.Println("deferred: third")
	fmt.Println("function body running")
	// Output order: body running, then third, second, first — LIFO.
}

func loopGotcha() {
	fmt.Println("--- loopGotcha (arguments evaluate at defer time) ---")
	for i := 0; i < 3; i++ {
		defer fmt.Println("gotcha defer:", i)
	}
	// Each defer's argument i is evaluated immediately as the loop runs,
	// so this schedules Println(0), Println(1), Println(2) in that
	// order — but LIFO means they print 2, 1, 0.
}

func loopFixed() {
	fmt.Println("--- loopFixed (explicit closure argument) ---")
	for i := 0; i < 3; i++ {
		defer func(n int) {
			fmt.Println("fixed defer:", n)
		}(i)
	}
	// Same LIFO order, same output as loopGotcha here — the fix matters
	// most when the deferred call closes over i without taking it as a
	// parameter, which reads the loop variable's final value instead.
}

func main() {
	orderedDefers()
	loopGotcha()
	loopFixed()
}
