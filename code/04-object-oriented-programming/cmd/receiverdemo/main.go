// Command receiverdemo shows the value-vs-pointer receiver decision that
// stands in for "instance methods" in Go.
//
// Item pairs data (the struct) with behavior (methods) — that's the
// whole "class" story in Go, no keyword required.
//
// Run:
//
//	go run ./cmd/receiverdemo
package main

import "fmt"

// Item is a plain struct: just data, nothing more.
type Item struct {
	Name     string
	Quantity int
	UnitCost float64
}

// TotalValue has a VALUE receiver — it operates on a copy of the Item.
// Barely different from a free function TotalValue(i Item) float64; the
// receiver syntax is sugar for writing it.TotalValue() instead of
// TotalValue(it).
func (i Item) TotalValue() float64 {
	return float64(i.Quantity) * i.UnitCost
}

func (i *Item) SellItem() {
	i.Quantity--
	fmt.Printf("Sold one %s, %d remaining", i.Name, i.Quantity)
}

// ApplyDiscount has a POINTER receiver — it operates on the original
// Item, so the mutation is visible to the caller after the call returns.
func (i *Item) ApplyDiscount(pct float64) {
	i.UnitCost *= (1 - pct/100)
}

func main() {

	// var it Item
	// it := Item{Name: "", Quantity: 0, UnitCost: 0}

	it := Item{Name: "Widget", Quantity: 10, UnitCost: 5}
	fmt.Printf("%s: qty=%d unitCost=%.2f total=%.2f\n",
		it.Name, it.Quantity, it.UnitCost, it.TotalValue())

	it.SellItem()
	fmt.Println(it)

	// Go auto-takes the address of it here: it.ApplyDiscount(10) is
	// shorthand for (&it).ApplyDiscount(10). Works because it is a
	// local, addressable variable.
	it.ApplyDiscount(10)
	fmt.Printf("after 10%% discount: unitCost=%.2f total=%.2f\n",
		it.UnitCost, it.TotalValue())

	// Prove the value receiver really does get a copy: TotalValue can't
	// have mutated it, so the discount above is the only change.
	before := it
	_ = before.TotalValue() // computed from a copy, changes nothing
	fmt.Printf("unitCost still %.2f after calling a value-receiver method\n", it.UnitCost)
}
