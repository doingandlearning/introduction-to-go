// Command inventory is the entry point for Lab 2: a small warehouse
// stock report exercising multiple return values, zero values, defer,
// a second package, and the core fmt verbs.
//
//	go run ./cmd/inventory
package main

import (
	"fmt"
	"io"
	"os"

	"example.com/core-lab/internal/mathutils"
)

// Item is one line of warehouse stock.
type Item struct {
	Name      string
	Quantity  int
	UnitPrice float64
}

func sampleItems() []Item {
	return []Item{
		{Name: "Bolts (box of 100)", Quantity: 40, UnitPrice: 4.50},
		{Name: "Hinges", Quantity: 120, UnitPrice: 1.25},
		{Name: "Paint (1L)", Quantity: 15, UnitPrice: 8.99},
	}
}

// totalQuantity sums Quantity across items using mathutils.Add.
//
// (Exercise 5): mathutils.Add is a TODO in mathutils.go — until you
// implement it, this will compile and run but always report 0, since
// the placeholder Add always returns 0.
func totalQuantity(items []Item) int {
	quantities := make([]int, 0, len(items))
	for _, it := range items {
		quantities = append(quantities, it.Quantity)
	}
	return mathutils.Add(quantities...)
}

// totalValue sums Quantity * UnitPrice across items.
func totalValue(items []Item) float64 {
	total := 0.0
	for _, it := range items {
		total += float64(it.Quantity) * it.UnitPrice
	}
	return total
}

// averageUnitPrice divides totalValue by totalQuantity using
// mathutils.SafeDivide, so a zero-quantity warehouse reports a message
// instead of crashing.
//
// (Exercise 1): mathutils.SafeDivide is a TODO in mathutils.go — until
// you implement it, this will compile and run but always report 0,
// since the placeholder SafeDivide always returns (0, nil).
func averageUnitPrice(value float64, qty int) float64 {
	avg, err := mathutils.SafeDivide(value, float64(qty))
	if err != nil {
		fmt.Println("could not compute average unit price:", err)
		return 0
	}
	return avg
}

// reportZeroValues declares one variable of several basic types, plus
// an unassigned Item, without initializing any of them, and prints each
// to w. A pre-written test (main_test.go) checks every label below is
// present — run `go test ./...` before you start, then again once
// you've filled this in.
//
// TODO (Exercise 2): declare and print unset int, float64, string,
// bool, *Item, and Item variables using %v (and %+v for the struct).
// Label each line the same way the test expects: "int:", "float64:",
// "string:", "bool:", "*Item:", "Item:".
func reportZeroValues(w io.Writer) {
	fmt.Fprintln(w, "--- zero values ---")
	// TODO: replace this placeholder.
}

// auditReport has three defers already written — nothing to implement.
// Read it, predict the output order, then run it (Exercise 3).
func auditReport() {
	fmt.Println("--- audit report ---")
	defer fmt.Println("audit: closing report file")
	defer fmt.Println("audit: flushing summary buffer")
	defer fmt.Println("audit: releasing warehouse lock")
	fmt.Println("audit: report body written")
}

// closeZonesBuggy reproduces the classic loop-variable defer gotcha:
// nothing to implement here either, it's meant to be run and observed.
func closeZonesBuggy() {
	fmt.Println("--- closing zones (buggy) ---")
	zones := []string{"A", "B", "C"}
	for _, z := range zones {
		defer fmt.Println("closing zone:", z)
	}
}

// closeZonesFixed should produce the same zone names as closeZonesBuggy,
// but each deferred call should capture the zone it was scheduled for.
// A pre-written test checks the three "closing zone: X" lines come out
// in LIFO order (C, then B, then A) — run `go test ./...` to see it
// fail first.
//
// TODO (Exercise 4): loop over the same zones slice, deferring a
// closure that takes the zone name as an explicit parameter, e.g.
// defer func(zone string) { fmt.Fprintln(w, "closing zone:", zone) }(z).
func closeZonesFixed(w io.Writer) {
	fmt.Fprintln(w, "--- closing zones (fixed) ---")
	// TODO: replace this placeholder.

}

// printItemFormats prints one Item with %v, %+v, %#v, and %T. A
// pre-written test checks all four verb labels and %T's package-qualified
// type name show up in the output.
//
// TODO (Exercise 6): print sample using all four verbs, one per line,
// each labeled with the verb that produced it — fmt.Fprintf(w, "%%v   ->
// %v\n", sample), and so on for %+v, %#v, %T.
func printItemFormats(w io.Writer, sample Item) {
	fmt.Fprintln(w, "--- format comparison ---")
	// TODO: replace this placeholder.
}

func main() {
	items := sampleItems()

	reportZeroValues(os.Stdout)
	fmt.Println()

	auditReport()
	fmt.Println()

	closeZonesBuggy()
	closeZonesFixed(os.Stdout)
	fmt.Println()

	qty := totalQuantity(items)
	value := totalValue(items)
	avg := averageUnitPrice(value, qty)
	fmt.Printf("total quantity: %d, total value: %.2f, average unit price: %.2f\n", qty, value, avg)
	fmt.Println()

	printItemFormats(os.Stdout, items[0])
}
