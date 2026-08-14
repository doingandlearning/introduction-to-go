// Command billingdemo exercises the billing package's exported
// functions so `go run` prints something concrete without needing a
// debugger. The interesting content for this topic lives in
// billing_test.go, not here.
package main

import (
	"fmt"

	"example.com/testing-go/internal/billing"
)

func main() {
	price := 120.0
	fmt.Printf("Price:              %.2f\n", price)
	fmt.Printf("Ten percent off:    %.2f\n", billing.TenPercentOff(price))
	fmt.Printf("Tier discount:      %.2f\n", billing.TierDiscount(price))
	fmt.Printf("Applied strategy:   %.2f\n", billing.ApplyDiscount(price, billing.TenPercentOff))
	fmt.Printf("120 / 4:            %.2f\n", billing.Divide(120, 4))
}
