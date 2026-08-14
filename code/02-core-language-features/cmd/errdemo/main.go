// Command errdemo shows Go's multiple-return-value approach to error
// handling: no exceptions, no try/catch, just a value the caller checks.
//
//	go run ./cmd/errdemo
package main

import (
	"errors"
	"fmt"
	"log"
)

// divide returns an error instead of panicking or throwing when b is zero.
// The caller decides what "failure" means for them.
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	// Happy path: check the error, it's nil, move on.
	result, err := divide(10, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("10 / 2 = %v\n", result)

	// Failure path: the same pattern, this time the error is real.
	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("10 / 0 failed: %v (result is the zero value: %v)\n", err, result)
	}

	// fmt.Errorf is the more common way to build an error, since it lets
	// you format context straight into the message.
	_, err = divide(1, 0)
	if err != nil {
		wrapped := fmt.Errorf("computing ratio: %w", err)
		fmt.Println("wrapped error:", wrapped)
	}
}
