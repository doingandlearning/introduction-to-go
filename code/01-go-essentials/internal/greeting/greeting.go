// Package greeting demonstrates Go's only visibility rule: capitalize a
// name to export it from the package, lowercase it to keep it private.
package greeting

import "fmt"

// Greet is exported — its name starts with an uppercase letter, so any
// package that imports "greeting" can call it.
func Greet(name string) string {
	return fmt.Sprintf("%s, %s", greetingPrefix(), name)
}

// greetingPrefix is unexported — lowercase, so it's only callable from
// inside this package. Nothing outside "greeting" can see it, and there's
// no keyword that would let you change that. The package boundary is the
// only boundary Go has.
func greetingPrefix() string {
	return "Hello"
}
