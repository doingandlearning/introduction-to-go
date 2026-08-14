// Command nilinterfacedemo reproduces the famous nil-pointer-in-an-
// interface gotcha, then shows the fix. This trips up experienced Go
// developers regularly, not just beginners - it's a genuinely sharp
// corner of the language, not a "you're not getting it" situation.
//
//	go run ./cmd/nilinterfacedemo
package main

import "fmt"

// MyError has a pointer receiver on Error, which matters here: only
// *MyError satisfies the error interface, not MyError by value.
type MyError struct{}

func (e *MyError) Error() string { return "something broke" }

// leaky reproduces the gotcha. It declares e as a concrete *MyError,
// leaves it nil, and returns it as an error. The returned interface
// value is the pair (*MyError, nil) - the type half is non-nil, so
// comparing the interface to nil is false, even though e "looks" nil.
func leaky() error {
	var e *MyError = nil
	// ... imagine real logic here that never assigns to e ...
	return e
}

// fixed avoids the bug by never letting a concrete pointer type reach
// the interface return slot in the success path - it returns a literal
// nil instead of a variable declared as *MyError.
func fixed() error {
	somethingWentWrong := false
	if somethingWentWrong {
		return &MyError{}
	}
	return nil
}

func main() {
	err := leaky()
	fmt.Println("leaky():  err == nil ->", err == nil) // false - the gotcha

	err = fixed()
	fmt.Println("fixed():  err == nil ->", err == nil) // true - correct
}
