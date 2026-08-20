package main

import "fmt"

func main() {
	a := [3]int{1, 2, 3}
	fmt.Println(len(a), cap(a))

	s := []int{1, 2, 3} // slice literal
	s = append(s, 4)    // grows as needed
	fmt.Println(len(s), cap(s))
	fmt.Println(s[3])
	fmt.Println(s[1:])
	fmt.Println(s[:2])
}
