// Package orders holds the pieces the Lab 6 pipeline is built from: a
// counter closure, two generic helpers, a functional-options order
// builder, and an invoice type used to demonstrate method values.
package orders

// ---------------------------------------------------------------------
// Exercise 1: NewOrderCounter
// ---------------------------------------------------------------------

// NewOrderCounter returns a closure that hands out order numbers
// starting at 1 and incrementing on every call. Each call to
// NewOrderCounter must produce a counter with its own independent
// starting state - two counters must never share a running total.
//
// TODO: implement this using a closure, the same shape as makeCounter
// from the lecture.
func NewOrderCounter() func() int {
	// TODO: replace this stub. It currently always returns 0, which is
	// wrong - every call should return a fresh closure over its own
	// count variable.
	return func() int {
		return 0
	}
}

// ---------------------------------------------------------------------
// Exercise 2: Filter
// ---------------------------------------------------------------------

// Filter returns the elements of items for which predicate returns
// true. It must work for any element type T without modification.
//
// TODO: implement this generically.
func Filter[T any](items []T, predicate func(T) bool) []T {
	// TODO: replace this stub. It currently returns every item
	// unfiltered.
	return items
}

// ---------------------------------------------------------------------
// Exercise 3: Map + Drink.Dollars
// ---------------------------------------------------------------------

// Drink is a menu item. PriceCents is the stored value; Dollars is a
// derived value computed from it.
type Drink struct {
	Name       string
	PriceCents int
}

// Dollars converts PriceCents to a float64 number of dollars.
//
// TODO: implement this. Dividing an int by 100 in Go uses integer
// division - you'll need to convert to float64 before dividing.
func (d Drink) Dollars() float64 {
	// TODO: replace this stub.
	return 0
}

// Map applies transform to every element of items and returns the
// results. It must work for any pair of types T, U without
// modification.
//
// TODO: implement this generically.
func Map[T, U any](items []T, transform func(T) U) []U {
	// TODO: replace this stub. It currently returns an empty slice
	// regardless of input.
	return nil
}

// ---------------------------------------------------------------------
// Exercise 4: CoffeeOrder + functional options
// ---------------------------------------------------------------------

// CoffeeOrder is the thing NewCoffeeOrder builds. Fields are
// unexported - the only way to configure one from another package is
// through the CoffeeOption functions below.
type CoffeeOrder struct {
	size      string // "small", "medium", "large"
	extraShot bool
	oatMilk   bool
}

// CoffeeOption mutates a *CoffeeOrder. Each With... function below
// returns one of these - think of it as a modifier slip handed to the
// order.
//
// TODO: define this type.
type CoffeeOption func(*CoffeeOrder)

// WithSize overrides the default size.
//
// TODO: implement this.
func WithSize(size string) CoffeeOption {
	// TODO: replace this stub. It currently does nothing.
	return func(o *CoffeeOrder) {}
}

// WithExtraShot adds an extra shot to the order.
//
// TODO: implement this.
func WithExtraShot() CoffeeOption {
	// TODO: replace this stub. It currently does nothing.
	return func(o *CoffeeOrder) {}
}

// WithOatMilk switches the order to oat milk.
//
// TODO: implement this.
func WithOatMilk() CoffeeOption {
	// TODO: replace this stub. It currently does nothing.
	return func(o *CoffeeOrder) {}
}

// NewCoffeeOrder builds a CoffeeOrder with sensible defaults - size
// "medium", no extra shot, regular (not oat) milk - then applies
// whichever options were passed in.
//
// TODO: implement this. Defaults must survive untouched when no
// options are passed at all.
func NewCoffeeOrder(opts ...CoffeeOption) *CoffeeOrder {
	// TODO: replace this stub with real defaults and option application.
	return &CoffeeOrder{}
}

// Size, ExtraShot, and OatMilk expose the unexported fields read-only,
// so cmd/pipeline can print them without reaching into the struct
// directly.
func (o CoffeeOrder) Size() string    { return o.size }
func (o CoffeeOrder) ExtraShot() bool { return o.extraShot }
func (o CoffeeOrder) OatMilk() bool   { return o.oatMilk }

// ---------------------------------------------------------------------
// Exercise 5: Invoice.Total as a method value
// ---------------------------------------------------------------------

// Invoice represents the billed total for an order.
type Invoice struct {
	Subtotal float64
	TaxRate  float64
}

// Total returns Subtotal plus tax.
//
// TODO: implement this.
func (i Invoice) Total() float64 {
	// TODO: replace this stub.
	return 0
}
