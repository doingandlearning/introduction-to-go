// Command nilmap demonstrates the nil-slice vs nil-map asymmetry: a nil
// slice is safe to read AND append to (it allocates on first use), but a
// nil map is safe to read from and PANICS the moment you write to it.
//
//	go run ./cmd/nilmap
package main

import "fmt"

func main() {
	fmt.Println("-- nil slice: read and append both work --")
	var s []int
	fmt.Println("len(s) before append:", len(s), "s == nil:", s == nil)
	s = append(s, 1) // safe: append allocates a backing array on first use
	fmt.Println("s after append(s, 1):", s)

	fmt.Println()
	fmt.Println("-- nil map: read is safe, write panics --")
	var m map[string]int
	v, ok := m["anything"] // safe: reading a nil map returns the zero value
	fmt.Println(`m["anything"] ->`, v, ok)

	fmt.Println("about to write to a nil map — this will panic:")
	m["x"] = 1 // panic: assignment to entry in nil map

	// Unreachable. To fix, initialize first:
	//   m := make(map[string]int)
	// or
	//   m := map[string]int{}
	fmt.Println("never reached")
}
