// Package mathutils provides small numeric helpers for the Lab 2
// inventory CLI: a safe division that reports failure as a value, and a
// variadic sum used to total item quantities.
package mathutils

// SafeDivide returns a/b, or an error if b is zero, instead of letting
// the caller panic on a divide-by-zero.
//
// TODO (Exercise 1): implement this. It should:
//  1. Add an `import "fmt"` line above `package mathutils` — this file
//     doesn't import it yet, since nothing here needs it until you do.
//  2. Return an error built with fmt.Errorf when b == 0.
//  3. Otherwise return a/b and a nil error.
func SafeDivide(a, b float64) (float64, error) {
	// TODO: replace this placeholder.
	return 0, nil
}

// Add sums an arbitrary number of ints.
//
// TODO (Exercise 5): implement this. It should:
//  1. Accept a variadic []int parameter.
//  2. Sum the values with a for-range loop and return the total.
func Add(nums ...int) int {
	// TODO: replace this placeholder.
	return 0
}
