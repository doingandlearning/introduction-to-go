// Package billing contains small pricing and division helpers used
// throughout the Topic 12 testing examples. The business logic here is
// deliberately simple - the point of this topic is the tests sitting
// next to it, not the pricing rules themselves.
package billing

import "fmt"

// TenPercentOff returns price with a flat 10% discount applied.
func TenPercentOff(price float64) float64 {
	return price * 0.9
}

// DiscountStrategy is a function type - the same Strategy-pattern shape
// from Topic 9, reused here so the tests below exercise code you
// already recognise rather than a brand new toy example.
type DiscountStrategy func(price float64) float64

// ApplyDiscount runs price through whichever strategy it's handed. It
// doesn't know or care which one - only that it matches the
// DiscountStrategy shape.
func ApplyDiscount(price float64, strategy DiscountStrategy) float64 {
	return strategy(price)
}

// TierDiscount picks a discount tier based on price. Four distinct
// branches make it a good candidate for a table-driven test, and later
// for a coverage exercise.
func TierDiscount(price float64) float64 {
	switch {
	case price <= 0:
		return 0
	case price >= 100:
		return price * 0.8 // 20% off for big orders
	case price >= 50:
		return price * 0.9 // 10% off for mid-size orders
	default:
		flat := price - 5 // flat fiver off small orders
		if flat < 0 {
			return 0
		}
		return flat
	}
}

// Divide returns a/b. Unlike most Go functions in this course, it
// panics instead of returning an error when b is zero - written this
// way on purpose, so this package has something worth testing for a
// panic.
func Divide(a, b float64) float64 {
	if b == 0 {
		panic(fmt.Sprintf("billing: cannot divide %v by zero", a))
	}
	return a / b
}
