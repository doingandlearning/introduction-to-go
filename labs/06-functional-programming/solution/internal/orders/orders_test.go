// This test file is already complete — you're not writing it. It's the
// specification for Exercises 1-5: run `go test ./...` now, before
// touching orders.go, and every test below fails. Implement each
// TODO in orders.go until its matching test passes. Writing a test
// like this yourself is Topic 12's job, not this one's.
package orders

import (
	"math"
	"reflect"
	"testing"
)

// TestNewOrderCounter confirms two counters built from separate calls
// to NewOrderCounter track independent state (Exercise 1).
func TestNewOrderCounter(t *testing.T) {
	till1 := NewOrderCounter()
	till2 := NewOrderCounter()

	if got := till1(); got != 1 {
		t.Errorf("till1() 1st call = %d, want 1", got)
	}
	if got := till1(); got != 2 {
		t.Errorf("till1() 2nd call = %d, want 2", got)
	}
	if got := till2(); got != 1 {
		t.Errorf("till2() 1st call = %d, want 1 (independent of till1)", got)
	}
	if got := till1(); got != 3 {
		t.Errorf("till1() 3rd call = %d, want 3", got)
	}
	if got := till2(); got != 2 {
		t.Errorf("till2() 2nd call = %d, want 2", got)
	}
}

// TestFilter confirms Filter works unmodified across two unrelated
// element types (Exercise 2).
func TestFilter(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8}
	evens := Filter(ints, func(n int) bool { return n%2 == 0 })
	wantEvens := []int{2, 4, 6, 8}
	if !reflect.DeepEqual(evens, wantEvens) {
		t.Errorf("Filter(ints, even) = %v, want %v", evens, wantEvens)
	}

	names := []string{"tea", "latte", "cola", "espresso", "chai"}
	long := Filter(names, func(s string) bool { return len(s) > 3 })
	wantLong := []string{"latte", "cola", "espresso", "chai"}
	if !reflect.DeepEqual(long, wantLong) {
		t.Errorf("Filter(names, len>3) = %v, want %v", long, wantLong)
	}
}

// TestDollars confirms Drink.Dollars converts cents to dollars using
// float division, not integer division (Exercise 3).
func TestDollars(t *testing.T) {
	d := Drink{Name: "Latte", PriceCents: 375}
	got := d.Dollars()
	want := 3.75
	if math.Abs(got-want) > 0.0001 {
		t.Errorf("Drink{PriceCents: 375}.Dollars() = %v, want %v", got, want)
	}
}

// TestMap confirms Map works with a method expression
// (Drink.Dollars) as its transform, turning a []Drink into a []float64
// (Exercise 3).
func TestMap(t *testing.T) {
	menu := []Drink{
		{Name: "Espresso", PriceCents: 250},
		{Name: "Latte", PriceCents: 375},
		{Name: "Cold Brew", PriceCents: 425},
	}

	prices := Map(menu, Drink.Dollars)
	want := []float64{2.5, 3.75, 4.25}
	if len(prices) != len(want) {
		t.Fatalf("Map(menu, Drink.Dollars) returned %d prices, want %d", len(prices), len(want))
	}
	for i := range want {
		if math.Abs(prices[i]-want[i]) > 0.0001 {
			t.Errorf("prices[%d] = %v, want %v", i, prices[i], want[i])
		}
	}
}

// TestNewCoffeeOrder confirms defaults survive untouched with zero
// options, that a single option only changes the field it's
// responsible for, and that stacked options all apply together
// (Exercise 4).
func TestNewCoffeeOrder(t *testing.T) {
	// Zero options: defaults must survive untouched.
	plain := NewCoffeeOrder()
	if plain.Size() != "medium" {
		t.Errorf("zero options: Size() = %q, want %q", plain.Size(), "medium")
	}
	if plain.ExtraShot() {
		t.Errorf("zero options: ExtraShot() = true, want false")
	}
	if plain.OatMilk() {
		t.Errorf("zero options: OatMilk() = true, want false")
	}

	// One option: only the overridden field should change.
	large := NewCoffeeOrder(WithSize("large"))
	if large.Size() != "large" {
		t.Errorf("one option: Size() = %q, want %q", large.Size(), "large")
	}
	if large.ExtraShot() {
		t.Errorf("one option: ExtraShot() = true, want false")
	}
	if large.OatMilk() {
		t.Errorf("one option: OatMilk() = true, want false")
	}

	// Multiple options stacked: every override should apply together.
	loaded := NewCoffeeOrder(WithSize("small"), WithExtraShot(), WithOatMilk())
	if loaded.Size() != "small" {
		t.Errorf("stacked options: Size() = %q, want %q", loaded.Size(), "small")
	}
	if !loaded.ExtraShot() {
		t.Errorf("stacked options: ExtraShot() = false, want true")
	}
	if !loaded.OatMilk() {
		t.Errorf("stacked options: OatMilk() = false, want true")
	}
}

// TestInvoiceTotal confirms Total computes Subtotal plus tax, and that
// assigning it to a variable as a method value - no parentheses, no
// call - still works with no arguments once called, because it
// captured its receiver at assignment time (Exercise 5).
func TestInvoiceTotal(t *testing.T) {
	inv := Invoice{Subtotal: 12.50, TaxRate: 0.08}
	want := 13.5

	if got := inv.Total(); math.Abs(got-want) > 0.0001 {
		t.Errorf("inv.Total() = %v, want %v", got, want)
	}

	getTotal := inv.Total // method value: no parens, not called yet
	if got := getTotal(); math.Abs(got-want) > 0.0001 {
		t.Errorf("getTotal() (method value, called later) = %v, want %v", got, want)
	}
}
