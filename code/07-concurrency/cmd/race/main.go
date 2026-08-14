// Command race is INTENTIONALLY BROKEN — run it live to observe the race
// detector output, don't leave it running in CI or import this package
// anywhere.
//
// Two goroutines both increment a shared `counter` variable a thousand
// times each, with no synchronization at all. counter++ is a
// read-modify-write, not an atomic operation, so the two goroutines can
// interleave and lose increments.
//
// Run it normally first — it usually "succeeds" with a wrong number:
//
//	go run ./cmd/race
//
// Expected output (varies run to run, rarely exactly 2000):
//
//	final counter value: 1847
//
// Now run it with the race detector:
//
//	go run -race ./cmd/race
//
// Expected output: a WARNING: DATA RACE report naming the two goroutines,
// the exact file/line of each conflicting access, and whether each was a
// read or a write, e.g.:
//
//	WARNING: DATA RACE
//	Write at 0x00c0000140a0 by goroutine 8:
//	  main.increment.func1()
//	      .../cmd/race/main.go:29 +0x3c
//	Previous write at 0x00c0000140a0 by goroutine 7:
//	  main.increment.func1()
//	      .../cmd/race/main.go:29 +0x3c
//
// See cmd/race/fixed_mutex and cmd/race/fixed_channel below in the source
// (kept as comments, not separate builds) for the two idiomatic fixes
// discussed in the slides: a sync.Mutex, and a redesign where each
// goroutine sends increments down a channel to a single owning goroutine.
package main

import (
	"fmt"
	"sync"
)

var counter int

func increment(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		counter++ // read-modify-write with no synchronization — the race
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go increment(&wg)
	}

	wg.Wait()
	fmt.Println("final counter value:", counter)
}

// --- Fix 1: sync.Mutex ---------------------------------------------------
//
// var mu sync.Mutex
// var counter int
//
// func increment(wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for i := 0; i < 1000; i++ {
// 		mu.Lock()
// 		counter++
// 		mu.Unlock()
// 	}
// }
//
// Simple, direct, and completely correct for small shared state like a
// single counter. "Share memory by communicating" is Go's default, not a
// law — a mutex is often the right tool for exactly this shape of problem.

// --- Fix 2: channel + single owning goroutine -----------------------------
//
// func main() {
// 	incs := make(chan int)
// 	done := make(chan int)
//
// 	// The owning goroutine is the ONLY thing that ever touches total.
// 	go func() {
// 		total := 0
// 		for range incs {
// 			total++
// 		}
// 		done <- total
// 	}()
//
// 	var wg sync.WaitGroup
// 	for i := 0; i < 2; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for i := 0; i < 1000; i++ {
// 				incs <- 1 // send an increment instead of touching a shared var
// 			}
// 		}()
// 	}
//
// 	wg.Wait()
// 	close(incs)
// 	fmt.Println("final counter value:", <-done)
// }
