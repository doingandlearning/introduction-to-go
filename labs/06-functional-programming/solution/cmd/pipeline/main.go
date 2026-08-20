// Command pipeline is the entry point for Lab 6, a small order pipeline
// for the "Go Roasters" coffee shop. Each section below corresponds to
// one exercise in labs/exercise.md.
//
//	go run ./cmd/pipeline
package main

import (
	"fmt"

	"example.com/goroasters/internal/orders"
)

func main() {
	exerciseOneCounters()
	exerciseTwoFilter()
	exerciseThreeMap()
	exerciseFourOptions()
	exerciseFiveMethodValue()
}

// Exercise 1: two independent order-number counters, one per till.
func exerciseOneCounters() {
	fmt.Println("--- Exercise 1: independent counters ---")

	till1 := orders.NewOrderCounter()
	till2 := orders.NewOrderCounter()

	fmt.Println("till1:", till1()) // 1
	fmt.Println("till1:", till1()) // 2
	fmt.Println("till2:", till2()) // 1 - independent of till1
	fmt.Println("till1:", till1()) // 3
	fmt.Println("till2:", till2()) // 2
	fmt.Println()
}

// Exercise 2: the same Filter function over two unrelated element types.
func exerciseTwoFilter() {
	fmt.Println("--- Exercise 2: Filter over two types ---")

	orderSizes := []int{1, 2, 3, 4, 5, 6, 7, 8}
	evens := orders.Filter(orderSizes, func(n int) bool { return n%2 == 0 })
	fmt.Println("even sizes:", evens) // [2 4 6 8]

	drinkNames := []string{"tea", "latte", "cola", "espresso", "chai"}
	longNames := orders.Filter(drinkNames, func(s string) bool { return len(s) > 3 })
	fmt.Println("long names:", longNames) // [latte cola espresso chai]
	fmt.Println()
}

// Exercise 3: Map turns a []Drink into a []float64 of dollar prices.
func exerciseThreeMap() {
	fmt.Println("--- Exercise 3: Map over Drinks ---")

	menu := []orders.Drink{
		{Name: "Espresso", PriceCents: 250},
		{Name: "Latte", PriceCents: 375},
		{Name: "Cold Brew", PriceCents: 425},
	}

	prices := orders.Map(menu, orders.Drink.Dollars)
	fmt.Println("prices:", prices) // [2.5 3.75 4.25]
	fmt.Println()
}

// Exercise 4: functional options - zero, one, and two stacked.
func exerciseFourOptions() {
	fmt.Println("--- Exercise 4: functional options ---")

	plain := orders.NewCoffeeOrder()
	fmt.Printf("plain:  size=%s extraShot=%v oatMilk=%v\n",
		plain.Size(), plain.ExtraShot(), plain.OatMilk()) // medium false false

	large := orders.NewCoffeeOrder(orders.WithSize("large"))
	fmt.Printf("large:  size=%s extraShot=%v oatMilk=%v\n",
		large.Size(), large.ExtraShot(), large.OatMilk()) // large false false

	loaded := orders.NewCoffeeOrder(orders.WithSize("large"), orders.WithExtraShot(), orders.WithOatMilk())
	fmt.Printf("loaded: size=%s extraShot=%v oatMilk=%v\n",
		loaded.Size(), loaded.ExtraShot(), loaded.OatMilk()) // large true true
	fmt.Println()
}

// Exercise 5: method values - a method assigned to a variable without
// being called still remembers its receiver.
func exerciseFiveMethodValue() {
	fmt.Println("--- Exercise 5: method value ---")

	inv := orders.Invoice{Subtotal: 12.50, TaxRate: 0.08}
	fmt.Printf("invoice: %+v\n", inv)

	getTotal := inv.Total             // method value - no call, no parens
	fmt.Println("total:", getTotal()) // 13.5, no argument needed
}
