// Command hello is the smallest useful Go program: one package, one
// entry point, no runtime required once it's built.
//
// Try it three ways from this directory:
//
//	go run main.go        # compile + run, throws away the binary
//	go build && ./hello    # produces a standalone binary, then run it
//	go install             # builds and places the binary in $HOME/go/bin
package main

import "fmt"

func main() {
	fmt.Println("Hello, Go")
}
