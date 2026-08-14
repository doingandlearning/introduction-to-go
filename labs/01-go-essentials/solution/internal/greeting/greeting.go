// Package greeting builds greeting strings for the Lab 1 CLI.
package greeting

import "strings"

// Greet returns a shouted greeting for name, e.g. "HELLO, ADA!".
func Greet(name string) string {
	return shout("Hello, " + name + "!")
}

// shout is unexported on purpose — see Exercise 2, step 4. Nothing
// outside this package can call it; main.go has to go through Greet.
func shout(s string) string {
	return strings.ToUpper(s)
}
