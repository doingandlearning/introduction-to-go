// Package mathutils provides small numeric helpers for the Lab 2
// inventory CLI: a safe division that reports failure as a value, and a
// variadic sum used to total item quantities.
package mathutils

import "fmt"

// SafeDivide returns a/b, or an error if b is zero, instead of letting
// the caller panic on a divide-by-zero.
func SafeDivide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide %v by zero", a)
	}
	return a / b, nil
}

// Add sums an arbitrary number of ints.
func Add(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
