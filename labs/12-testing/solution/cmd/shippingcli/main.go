// Command shippingcli prints the shipping package's behaviour so you
// can see it work outside of a test run.
package main

import (
	"fmt"

	"example.com/testing-lab/internal/shipping"
)

func main() {
	weight := 12.0
	fmt.Printf("Weight (kg):        %.2f\n", weight)
	fmt.Printf("Zone rate:          %.2f\n", shipping.ZoneRate(weight))

	cost := 80.0
	fmt.Printf("Cost:               %.2f\n", cost)
	fmt.Printf("Express surcharge:  %.2f\n", shipping.ExpressSurcharge(cost))
	fmt.Printf("Applied strategy:   %.2f\n", shipping.ApplyRate(cost, shipping.ExpressSurcharge))
	fmt.Printf("Cost per parcel:    %.2f\n", shipping.CostPerParcel(90, 3))
}
