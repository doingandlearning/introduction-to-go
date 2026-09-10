package main

import "fmt"

// Strategy interface
type ShippingStrategy interface {
	Calculate(weightKg float64) float64
}

// Concrete strategies
type StandardShipping struct{}

func (s StandardShipping) Calculate(weightKg float64) float64 {
	return weightKg * 2.50
}

type ExpressShipping struct{}

func (s ExpressShipping) Calculate(weightKg float64) float64 {
	return weightKg*4.00 + 10 // flat surcharge
}

type FreeShipping struct{}

func (s FreeShipping) Calculate(weightKg float64) float64 {
	return 0
}

// Context that uses a strategy
type Order struct {
	WeightKg float64
	Strategy ShippingStrategy
}

func (o Order) ShippingCost() float64 {
	return o.Strategy.Calculate(o.WeightKg)
}

func main() {
	order := Order{WeightKg: 3.2, Strategy: ExpressShipping{}}
	fmt.Printf("Cost: £%.2f\n", order.ShippingCost())

	order.Strategy = FreeShipping{}
	fmt.Printf("Cost: £%.2f\n", order.ShippingCost())
}
