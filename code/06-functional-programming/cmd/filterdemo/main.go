// Command filterdemo shows one generic Filter function working across
// two unrelated element types, and a generic Map function turning a
// slice of structs into a slice of plain float64s using a method value.
//
//	go run ./cmd/filterdemo
package main

import "fmt"

// Filter returns the elements of items for which predicate returns true.
// T any means this compiles once and works for every element type - no
// interface{} and no type assertions anywhere in this function.
func Filter[T any](items []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map turns a []T into a []U by applying transform to every element.
func Map[T, U any](items []T, transform func(T) U) []U {
	result := make([]U, 0, len(items))
	for _, item := range items {
		result = append(result, transform(item))
	}
	return result
}

// Drink is the struct Map will transform below. Dollars is a derived
// value, computed from PriceCents rather than stored directly.
type Drink struct {
	Name       string
	PriceCents int
}

func (d Drink) Dollars() float64 {
	return float64(d.PriceCents) / 100
}

func main() {
	// Filter over ints.
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Println("evens:", evens)

	// Same Filter function, completely different element type.
	languages := []string{"go", "rust", "c", "python", "zig"}
	long := Filter(languages, func(s string) bool { return len(s) > 3 })
	fmt.Println("long names:", long)

	// Map: []Drink -> []float64. Drink.Dollars here is a *method
	// expression* (Type.Method), which Go turns into an ordinary
	// func(Drink) float64 - exactly the shape Map's transform parameter
	// wants. Contrast this with a *method value* (menu[0].Dollars),
	// which binds a specific receiver instead of taking one as an
	// argument - see the lab for that version.
	menu := []Drink{
		{Name: "Espresso", PriceCents: 250},
		{Name: "Latte", PriceCents: 375},
		{Name: "Cold Brew", PriceCents: 425},
	}
	prices := Map(menu, Drink.Dollars)
	fmt.Println("prices:", prices)
}
