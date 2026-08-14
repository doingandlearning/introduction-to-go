// Command singletonlab exercises the registry Singleton across
// Exercises 1-3. See labs/exercise.md for the tasks.
//
//	go run ./cmd/singletonlab
package main

import (
	"fmt"

	"example.com/patterns-lab/internal/registry"
)

func main() {
	// Exercise 1: confirm two calls return the same pointer.
	a := registry.GetSettings()
	b := registry.GetSettings()
	fmt.Printf("a=%p b=%p same=%v\n", a, b, a == b)

	// Exercise 2 & 3: uncomment this block once instructed to in the lab.
	// It fires many goroutines at GetSettings concurrently.
	//
	// const goroutines = 50
	// var wg sync.WaitGroup
	// wg.Add(goroutines)
	// for i := 0; i < goroutines; i++ {
	// 	go func() {
	// 		defer wg.Done()
	// 		_ = registry.GetSettings()
	// 	}()
	// }
	// wg.Wait()
	// fmt.Println("all goroutines finished")
}
