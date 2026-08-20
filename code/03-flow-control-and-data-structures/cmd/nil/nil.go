package main

import "fmt"

func main() {
	var test []int
	test = append(test, 3)
	fmt.Println(test)

	var test2 map[string]bool //= make(map[string]bool)
	test2["name"] = false
	fmt.Println(test2)
}
