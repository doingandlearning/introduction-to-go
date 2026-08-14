// Command zerodemo declares one variable of several basic types without
// assigning it, then prints each one to show Go's zero values in place
// of null/undefined/None.
//
//	go run ./cmd/zerodemo
package main

import "fmt"

type Point struct {
	X, Y int
}

func main() {
	var i int
	var f float64
	var s string
	var b bool
	var p *int
	var sl []int
	var m map[string]int
	var pt Point

	fmt.Printf("int:        %v\n", i)
	fmt.Printf("float64:    %v\n", f)
	fmt.Printf("string:     %q\n", s)
	fmt.Printf("bool:       %v\n", b)
	fmt.Printf("*int:       %v\n", p)
	fmt.Printf("[]int:      %v (nil? %t)\n", sl, sl == nil)
	fmt.Printf("map:        %v (nil? %t)\n", m, m == nil)
	fmt.Printf("Point:      %+v (usable with no constructor)\n", pt)
}
