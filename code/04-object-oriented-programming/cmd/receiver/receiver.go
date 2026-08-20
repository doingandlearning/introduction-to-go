package main

import "fmt"

type Vehicle struct {
	Plate   string
	Mileage int
}

func (v Vehicle) EstServiceCost() float64 {
	return float64(v.Mileage) * 0.05
}

func EstimateServiceCost(v Vehicle) float64 {
	return float64(v.Mileage) * 0.05
}

func main() {
	car := Vehicle{
		Plate:   "ABC123",
		Mileage: 10000,
	}

	fmt.Println(car.EstServiceCost())
	fmt.Println(EstimateServiceCost(car))
}
