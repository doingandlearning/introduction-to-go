// This file is a REFERENCE for Exercise 6, task 3 — it is not part of the
// buildable `standalone/race` program (it lives in package main alongside
// main.go, which already defines main() and deliveryCount, so this file
// is kept excluded from the build via the "ignore" build tag below).
//
// To try it: copy this file's contents into its own directory with its
// own go.mod (or comment out main.go temporarily), then run:
//
//	go run -race .
//
// It should report no race.
//
//go:build ignore

package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var deliveryCountMutex int

func recordDeliveriesMutex(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		mu.Lock()
		deliveryCountMutex++
		mu.Unlock()
	}
}

func mainMutexFix() {
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go recordDeliveriesMutex(&wg)
	}
	wg.Wait()
	fmt.Println("final delivery count:", deliveryCountMutex)
}
