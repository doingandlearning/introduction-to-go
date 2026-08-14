// Command vetdemo compiles cleanly but has a bug the compiler doesn't
// catch: a Printf format string with the wrong number of arguments.
//
// Run:
//
//	go build ./...      # succeeds — this is valid Go
//	go vet ./...         # flags the Printf/argument mismatch below
//
// This is the layer distinction worth internalizing early: the compiler
// checks that your program is legal Go. go vet checks that it's probably
// not a mistake. They catch different things.
package main

import "fmt"

func main() {
	name := "Go"
	// Two verbs, one argument. Compiles fine. go vet will not let it slide.
	fmt.Printf("Hello, %s! You are %d years old.\n", name)
}
