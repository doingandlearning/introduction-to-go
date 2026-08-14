// Command itemctor shows Go's constructor convention: there is no
// constructor syntax, just a plain function that builds and validates a
// struct. By convention it's named NewX.
//
// Run:
//
//	go run ./cmd/itemctor
package main

import (
	"errors"
	"fmt"
)

// Item is the same shape as in receiverdemo — data only.
type Item struct {
	Name     string
	Quantity int
	UnitCost float64
}

// NewItem is the conventional constructor: it validates inputs and
// returns a pointer plus an error, the same two-value pattern used
// everywhere else in Go for "this might fail."
func NewItem(name string, qty int, cost float64) (*Item, error) {
	if qty < 0 {
		return nil, fmt.Errorf("quantity cannot be negative, got %d", qty)
	}
	if cost < 0 {
		return nil, errors.New("unit cost cannot be negative")
	}
	return &Item{Name: name, Quantity: qty, UnitCost: cost}, nil
}

func main() {
	// Success path.
	it, err := NewItem("Widget", 10, 5.00)
	if err != nil {
		fmt.Println("unexpected error:", err)
	} else {
		fmt.Printf("created: %+v\n", *it)
	}

	// Failure path — NewItem rejects it, main handles the error instead
	// of crashing or silently constructing a broken Item.
	bad, err := NewItem("Gadget", -3, 5.00)
	if err != nil {
		fmt.Println("rejected as expected:", err)
	} else {
		fmt.Printf("should not have succeeded: %+v\n", *bad)
	}

	// Reminder: nothing stops you from skipping NewItem entirely. The
	// zero value is a valid Item — see the OOP slides for why that
	// matters coming from Java, where an unconstructed object doesn't
	// exist at all.
	var zero Item
	fmt.Printf("zero value: %+v (still a usable Item)\n", zero)
}
