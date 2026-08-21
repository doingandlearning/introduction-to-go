package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("main.go")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fmt.Println("Closing file")
		f.Close() // Adds to a defer stack.
	}() // IIFE - Immediately Invoked Function Expression - runs once the defer is taken from the stack.

	tmp, err := os.CreateTemp("", "tmpfile")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name()) // Clean up the temp file when we're done.
	defer tmp.Close()           // Adds to a defer stack.

	// do something with f
	fmt.Println(f)
	fmt.Println(incrementer()) // Adds to a defer stack.
}

func incrementer() string {
	for i := 0; i < 10; i++ {
		defer fmt.Println(i) // Adds to a defer stack.
	}
	return "All done"
}
