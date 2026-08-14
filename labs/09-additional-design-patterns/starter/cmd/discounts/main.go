// Command discounts is Exercise 1: build the same three ticket discounts
// two ways - as a function type, then as an interface + structs - and
// compare the ceremony.
package main

import "fmt"

// --- Part A: function-type Strategy ---

// TODO (Exercise 1a): declare a TicketDiscount function type - a
// func(price float64) float64, exactly like the lecture's
// DiscountStrategy.

// TODO (Exercise 1a): implement three functions matching that type:
//   NoDiscount        - returns price unchanged
//   MemberDiscount     - returns price * 0.85 (15% off)
//   EarlyBirdDiscount  - returns price * 0.70 (30% off)

// TODO (Exercise 1a): implement ApplyTicketDiscount(price float64,
// d TicketDiscount) float64 that just calls d(price).

// --- Part B: interface-based Strategy ---

// TODO (Exercise 1b): declare a Discounter interface with one method,
// Apply(price float64) float64.

// TODO (Exercise 1b): implement three structs - noDiscount,
// memberDiscount, earlyBirdDiscount - each with an Apply method
// matching the same three discounts as Part A.

func main() {
	price := 200.0

	// TODO (Exercise 1a): call ApplyTicketDiscount with each of the three
	// functions and print the results.

	fmt.Println("price before any discount:", price)

	// TODO (Exercise 1b): construct each of the three Discounter structs,
	// call Apply on each, and print the results.

	// Once both parts compile and print matching numbers, count the
	// lines each version needed. Which one would you rather add a fourth
	// discount to?
}
