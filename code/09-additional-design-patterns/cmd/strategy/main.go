// Command strategy demonstrates the Strategy pattern the Go-native way:
// a function type, not an interface + a family of classes.
package main

import "fmt"

// DiscountStrategy is a function type - Go's version of the "algorithm
// slot" that Java would fill with a Strategy interface plus one class
// per implementation.
type DiscountStrategy func(price float64) float64

func NoDiscount(price float64) float64 {
	return price
}

func TenPercentOff(price float64) float64 {
	return price * 0.9
}

func TwentyPercentOff(price float64) float64 {
	return price * 0.8
}

// ApplyDiscount doesn't know or care which strategy it was handed - it
// only knows the shape: func(float64) float64.
func ApplyDiscount(price float64, strategy DiscountStrategy) float64 {
	return strategy(price)
}

func main() {
	price := 100.0

	fmt.Printf("No discount:     %.2f\n", ApplyDiscount(price, NoDiscount))
	fmt.Printf("10%% off:         %.2f\n", ApplyDiscount(price, TenPercentOff))
	fmt.Printf("20%% off:         %.2f\n", ApplyDiscount(price, TwentyPercentOff))

	// Strategies are ordinary values - store them in a map, pick one at
	// runtime, no factory or registry class required.
	strategies := map[string]DiscountStrategy{
		"none":   NoDiscount,
		"ten":    TenPercentOff,
		"twenty": TwentyPercentOff,
	}

	choice := "twenty"
	fmt.Printf("Chosen strategy %q: %.2f\n", choice, ApplyDiscount(price, strategies[choice]))
}
