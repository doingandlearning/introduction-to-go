// Command mapgotcha demonstrates a real Go gotcha: map values are not
// addressable, so you cannot call a pointer-receiver method directly on
// a value stored in a map.
//
// Run as-is:
//
//	go run ./cmd/mapgotcha
//
// To see the compile error live, uncomment the line marked "DOES NOT
// COMPILE" below and run `go build ./cmd/mapgotcha` again. Read the
// error, then re-comment it — the working fix follows right after it.
package main

import "fmt"

type Item struct {
	Name     string
	Quantity int
	UnitCost float64
}

func (i *Item) ApplyDiscount(pct float64) {
	i.UnitCost *= (1 - pct/100)
}

func main() {
	catalog := map[string]Item{
		"widget": {Name: "Widget", Quantity: 10, UnitCost: 5},
		"gadget": {Name: "Gadget", Quantity: 4, UnitCost: 20},
	}

	// DOES NOT COMPILE — uncomment to trigger it live:
	// catalog["widget"].ApplyDiscount(10)
	//
	// error: cannot call pointer method on catalog["widget"]
	//        cannot take address of catalog["widget"]
	//
	// Why: Go can auto-take the address of a local variable (see
	// receiverdemo), but a map value has no stable memory address —
	// the map can rehash and move it at any time. Go refuses to let
	// you take that address at all, rather than hand you one that
	// might already be stale.

	// The fix: read the value out (copies it), mutate the local copy,
	// write the whole thing back into the map.
	widget := catalog["widget"]
	widget.ApplyDiscount(10)
	catalog["widget"] = widget

	fmt.Printf("widget unit cost after discount: %.2f\n", catalog["widget"].UnitCost)
	fmt.Printf("gadget untouched: %.2f\n", catalog["gadget"].UnitCost)
}
