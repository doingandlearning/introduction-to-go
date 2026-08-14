// Command discounts is Exercise 1: build the same three ticket discounts
// two ways - as a function type, then as an interface + structs - and
// compare the ceremony.
package main

import "fmt"

// --- Part A: function-type Strategy ---

type TicketDiscount func(price float64) float64

func NoDiscount(price float64) float64       { return price }
func MemberDiscount(price float64) float64   { return price * 0.85 }
func EarlyBirdDiscount(price float64) float64 { return price * 0.70 }

func ApplyTicketDiscount(price float64, d TicketDiscount) float64 {
	return d(price)
}

// --- Part B: interface-based Strategy ---

type Discounter interface {
	Apply(price float64) float64
}

type noDiscount struct{}

func (noDiscount) Apply(price float64) float64 { return price }

type memberDiscount struct{}

func (memberDiscount) Apply(price float64) float64 { return price * 0.85 }

type earlyBirdDiscount struct{}

func (earlyBirdDiscount) Apply(price float64) float64 { return price * 0.70 }

func main() {
	price := 200.0
	fmt.Println("price before any discount:", price)

	fmt.Println("-- Part A: function type --")
	fmt.Printf("no discount:    %.2f\n", ApplyTicketDiscount(price, NoDiscount))
	fmt.Printf("member:         %.2f\n", ApplyTicketDiscount(price, MemberDiscount))
	fmt.Printf("early bird:     %.2f\n", ApplyTicketDiscount(price, EarlyBirdDiscount))

	fmt.Println("-- Part B: interface + structs --")
	discounters := []Discounter{noDiscount{}, memberDiscount{}, earlyBirdDiscount{}}
	for _, d := range discounters {
		fmt.Printf("%.2f\n", d.Apply(price))
	}

	// Part A: 1 type declaration + 3 one-line functions + 1 apply function.
	// Part B: 1 interface + 3 structs + 3 methods - roughly double the
	// declarations for the same three behaviors, with zero extra
	// capability, because none of these discounts needs more than one
	// method or hidden state.
}
