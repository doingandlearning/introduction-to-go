// Command greeter is the entry point for Lab 1. It expects a name on the
// command line and prints a greeting built by the internal/greeting
// package.
//
//	go run ./cmd/greeter Ada
package main

import (
	"fmt"
	"os"

	"example.com/greeter/internal/greeting"
)

func main() {
	name := "World"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	fmt.Println(greeting.Greet(name))

	// Exercise 3: call farewell(name) here once you've added it in
	// farewell.go, and print its result too.
}
