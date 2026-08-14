// This file is a REFERENCE for Exercise 6, task 4 — same build-tag setup
// as main_mutex.go: excluded from the default build so main.go's package
// still compiles standalone. Copy it out to try it for real.
//
//	go run -race .
//
// It should report no race — no shared variable is ever touched by more
// than one goroutine here; every increment travels down a channel to a
// single owning goroutine instead.
//
//go:build ignore

package main

import (
	"fmt"
	"sync"
)

func mainChannelFix() {
	incs := make(chan int)
	done := make(chan int)

	// The owning goroutine is the ONLY thing that ever touches total.
	go func() {
		total := 0
		for range incs {
			total++
		}
		done <- total
	}()

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				incs <- 1 // send an increment instead of sharing a variable
			}
		}()
	}

	wg.Wait()
	close(incs)
	fmt.Println("final delivery count:", <-done)
}
