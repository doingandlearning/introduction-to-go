package main

import "fmt"

// var ageIntTwo int = 5

// func test() {

// 	ageInt := float32(5) // type is inferred!

// 	y := float32(1.1)

// 	z := ageInt + y
// 	x++
// }

func main() {

	type Status float32

	const (
		StatusPending Status = iota * .25 // 1
		StatusActive                      // 2
		StatusDone                        // 3
	)

	fmt.Printf("%v", StatusDone)
}
