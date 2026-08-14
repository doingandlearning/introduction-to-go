// Command deadlock is INTENTIONALLY BROKEN — run it live to observe the
// crash, don't leave it running in CI or import this package anywhere.
//
// main sends on an unbuffered channel that nobody ever receives from.
// Because every goroutine in the program (here, just main) is asleep
// waiting on something that will never happen, the Go runtime detects the
// whole-program deadlock and crashes loudly instead of hanging silently:
//
//	go run ./cmd/deadlock
//
// Expected output:
//
//	fatal error: all goroutines are asleep - deadlock!
//
//	goroutine 1 [chan send]:
//	main.main()
//		.../cmd/deadlock/main.go:20 +0x25
//	exit status 2
//
// Note what this detector does NOT catch: if a THIRD goroutine were still
// busy doing something else while two other goroutines deadlocked on each
// other, the runtime would see the program as "alive" and report nothing.
// That version of a deadlock is silent and much harder to find — see the
// speaker notes on the "deadlock the runtime won't catch you" slide.
package main

func main() {
	ch := make(chan int) // unbuffered

	ch <- 42 // nobody is receiving — blocks forever, deadlock
}
