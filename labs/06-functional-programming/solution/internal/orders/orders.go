// Package orders holds the pieces the Lab 6 pipeline is built from: a
// counter closure, two generic helpers, a functional-options order
// builder, and an invoice type used to demonstrate method values.
package orders

// ---------------------------------------------------------------------
// Exercise 1: NewOrderCounter
// ---------------------------------------------------------------------

// NewOrderCounter returns a closure that hands out order numbers
// starting at 1 and incrementing on every call. Each call to
// NewOrderCounter produces a counter with its own independent count -
// two counters never share state, because each closure captures its
// own count variable.
func NewOrderCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// ---------------------------------------------------------------------
// Exercise 2: Filter
// ---------------------------------------------------------------------

// Filter returns the elements of items for which predicate returns
// true. Works for any element type T.
func Filter[T any](items []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
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
func (d Drink) Dollars() float64 {
	return float64(d.PriceCents) / 100
}

// Map applies transform to every element of items and returns the
// results. Works for any pair of types T, U.
func Map[T, U any](items []T, transform func(T) U) []U {
	result := make([]U, 0, len(items))
	for _, item := range items {
		result = append(result, transform(item))
	}
	return result
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

// CoffeeOption mutates a *CoffeeOrder - a modifier slip handed to the
// order.
type CoffeeOption func(*CoffeeOrder)

// WithSize overrides the default size.
func WithSize(size string) CoffeeOption {
	return func(o *CoffeeOrder) { o.size = size }
}

// WithExtraShot adds an extra shot to the order.
func WithExtraShot() CoffeeOption {
	return func(o *CoffeeOrder) { o.extraShot = true }
}

// WithOatMilk switches the order to oat milk.
func WithOatMilk() CoffeeOption {
	return func(o *CoffeeOrder) { o.oatMilk = true }
}

// NewCoffeeOrder builds a CoffeeOrder with sensible defaults - size
// "medium", no extra shot, regular (not oat) milk - then applies
// whichever options were passed in, in order. Defaults survive
// untouched when no options are passed at all.
func NewCoffeeOrder(opts ...CoffeeOption) *CoffeeOrder {
	o := &CoffeeOrder{
		size:      "medium",
		extraShot: false,
		oatMilk:   false,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
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
func (i Invoice) Total() float64 {
	return i.Subtotal * (1 + i.TaxRate)
}
