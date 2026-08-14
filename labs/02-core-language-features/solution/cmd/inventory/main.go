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
// to w.
func reportZeroValues(w io.Writer) {
	fmt.Fprintln(w, "--- zero values ---")

	var qty int
	var price float64
	var name string
	var inStock bool
	var ref *Item
	var blank Item

	fmt.Fprintf(w, "int:      %v\n", qty)
	fmt.Fprintf(w, "float64:  %v\n", price)
	fmt.Fprintf(w, "string:   %q\n", name)
	fmt.Fprintf(w, "bool:     %v\n", inStock)
	fmt.Fprintf(w, "*Item:    %v\n", ref)
	fmt.Fprintf(w, "Item:     %+v (usable with no constructor)\n", blank)
}

// auditReport has three defers already written — nothing to implement.
func auditReport() {
	fmt.Println("--- audit report ---")
	defer fmt.Println("audit: closing report file")
	defer fmt.Println("audit: flushing summary buffer")
	defer fmt.Println("audit: releasing warehouse lock")
	fmt.Println("audit: report body written")
	// Output order: body written, then lock, buffer, file — LIFO.
}

// closeZonesBuggy reproduces the classic loop-variable defer gotcha.
func closeZonesBuggy() {
	fmt.Println("--- closing zones (buggy) ---")
	zones := []string{"A", "B", "C"}
	for _, z := range zones {
		defer fmt.Println("closing zone:", z)
	}
	// Each defer's argument z is evaluated at defer time (A, then B,
	// then C), so LIFO prints C, B, A — arguably "correct" here since
	// the argument is captured by value. The real gotcha bites when the
	// deferred call closes over the loop variable instead of taking it
	// as a parameter; closeZonesFixed shows the safer pattern either way.
}

// closeZonesFixed passes the loop variable as an explicit parameter to
// the deferred closure, which is the idiomatic, gotcha-proof pattern.
func closeZonesFixed(w io.Writer) {
	fmt.Fprintln(w, "--- closing zones (fixed) ---")
	zones := []string{"A", "B", "C"}
	for _, z := range zones {
		defer func(zone string) {
			fmt.Fprintln(w, "closing zone:", zone)
		}(z)
	}
}

// printItemFormats prints one Item with %v, %+v, %#v, and %T.
func printItemFormats(w io.Writer, sample Item) {
	fmt.Fprintln(w, "--- format comparison ---")
	fmt.Fprintf(w, "%%v   -> %v\n", sample)
	fmt.Fprintf(w, "%%+v  -> %+v\n", sample)
	fmt.Fprintf(w, "%%#v  -> %#v\n", sample)
	fmt.Fprintf(w, "%%T   -> %T\n", sample)
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
