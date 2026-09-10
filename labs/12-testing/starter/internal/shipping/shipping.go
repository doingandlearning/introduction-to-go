// Package shipping contains small pricing and division helpers for the
// Topic 12 lab. Unlike every earlier lab in this course, the
// production code here is already finished — your job in this lab is
// to write the tests that sit next to it, not to complete TODOs in
// this file.
package shipping

import "fmt"

// ExpressSurcharge adds a flat 15% express-shipping surcharge to cost.
func ExpressSurcharge(cost float64) float64 {
	return cost * 1.15
}

// RateStrategy is a function type - the same Strategy-pattern shape
// from Topic 9.
type RateStrategy func(cost float64) float64

// ApplyRate runs cost through whichever strategy it's handed.
func ApplyRate(cost float64, strategy RateStrategy) float64 {
	return strategy(cost)
}

// ZoneRate picks a shipping rate based on parcel weight in kg:
//   - weightKg <= 0          -> 0
//   - weightKg >= 20          -> 4/kg heavy parcel rate
//   - 5 <= weightKg < 20      -> 5/kg standard parcel rate
//   - 0 < weightKg < 5        -> flat handling fee (10/kg minus an 8 rebate), floored at 0
func ZoneRate(weightKg float64) float64 {
	switch {
	case weightKg <= 0:
		return 0
	case weightKg >= 20:
		return weightKg * 4
	case weightKg >= 5:
		return weightKg * 5
	default:
		flat := weightKg*10 - 8
		if flat < 0 {
			return 0
		}
		return flat
	}
}

// CostPerParcel splits total across the given number of parcels. It
// panics instead of returning an error when parcels is zero - written
// this way on purpose, so this package has something worth testing for
// a panic.
func CostPerParcel(total, parcels float64) float64 {
	if parcels == 0 {
		panic(fmt.Sprintf("shipping: cannot split %v across zero parcels", total))
	}
	return total / parcels
}
