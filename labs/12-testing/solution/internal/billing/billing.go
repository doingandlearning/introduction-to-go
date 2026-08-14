// Package billing contains small pricing and division helpers for the
// Topic 12 lab. The production code is identical to starter/ - this
// lab is about the tests, not this file.
package billing

import "fmt"

// TenPercentOff returns price with a flat 10% discount applied.
func TenPercentOff(price float64) float64 {
	return price * 0.9
}

// DiscountStrategy is a function type - the same Strategy-pattern shape
// from Topic 9.
type DiscountStrategy func(price float64) float64

// ApplyDiscount runs price through whichever strategy it's handed.
func ApplyDiscount(price float64, strategy DiscountStrategy) float64 {
	return strategy(price)
}

// TierDiscount picks a discount tier based on price:
//   - price <= 0            -> 0
//   - price >= 100           -> 20% off
//   - 50 <= price < 100      -> 10% off
//   - 0 < price < 50         -> flat £5 off, floored at 0
func TierDiscount(price float64) float64 {
	switch {
	case price <= 0:
		return 0
	case price >= 100:
		return price * 0.8
	case price >= 50:
		return price * 0.9
	default:
		flat := price - 5
		if flat < 0 {
			return 0
		}
		return flat
	}
}

// Divide returns a/b. It panics instead of returning an error when b
// is zero - written this way on purpose, so this package has something
// worth testing for a panic.
func Divide(a, b float64) float64 {
	if b == 0 {
		panic(fmt.Sprintf("billing: cannot divide %v by zero", a))
	}
	return a / b
}
