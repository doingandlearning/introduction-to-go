// Command billingcli prints the billing package's behaviour so you can
// see it work outside of a test run.
package main

import (
	"fmt"

	"example.com/testing-lab/internal/billing"
)

func main() {
	price := 120.0
	fmt.Printf("Price:              %.2f\n", price)
	fmt.Printf("Ten percent off:    %.2f\n", billing.TenPercentOff(price))
	fmt.Printf("Tier discount:      %.2f\n", billing.TierDiscount(price))
	fmt.Printf("Applied strategy:   %.2f\n", billing.ApplyDiscount(price, billing.TenPercentOff))
	fmt.Printf("120 / 4:            %.2f\n", billing.Divide(120, 4))
}
