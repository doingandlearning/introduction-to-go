// Command closures demonstrates that a Go closure captures its own,
// independent copy of the enclosing-scope variables — not a shared
// global. Run it and read the interleaved output.
//
//	go run ./cmd/closures
package main

import "fmt"

// makeCounter returns a function that increments and returns its own
// private counter every time it's called. count only exists inside the
// closure returned here — nothing outside can reach it directly.
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	till1 := makeCounter()
	till2 := makeCounter()

	fmt.Println("till1:", till1()) // 1
	fmt.Println("till1:", till1()) // 2
	fmt.Println("till2:", till2()) // 1 - unaffected by till1's calls
	fmt.Println("till1:", till1()) // 3
	fmt.Println("till2:", till2()) // 2

	fmt.Println("\nFinal till1 count:", till1()) // 4
	fmt.Println("Final till2 count:", till2())   // 3
}
