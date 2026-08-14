// Command singletondemo confirms two things about the config Singleton:
//
//  1. Two separate calls to config.GetConfig() return the exact same
//     pointer.
//  2. Even when many goroutines call it concurrently, the underlying
//     loader only runs once.
//
//	go run ./cmd/singletondemo
package main

import (
	"fmt"
	"sync"

	"example.com/patterns-intro/internal/config"
)

func main() {
	// Part 1: same pointer from two call sites.
	a := config.GetConfig()
	b := config.GetConfig()
	fmt.Printf("same pointer? %v (a=%p b=%p)\n", a == b, a, b)

	// Part 2: concurrent access. Fire off many goroutines that all call
	// GetConfig at (roughly) the same time.
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = config.GetConfig()
		}()
	}
	wg.Wait()

	fmt.Printf("loadFromEnv ran %d time(s) across %d concurrent callers\n",
		config.LoadCount(), goroutines)
}
