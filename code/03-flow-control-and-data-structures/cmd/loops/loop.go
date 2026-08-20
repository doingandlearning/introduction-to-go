package main

import "fmt"

func main() {
	// while-style — a bare condition, no init or post statement
	count := 0
	for count < 5 {
		fmt.Println(count)
		count++
	}

	// infinite, exit via break
	books := []string{"Dune", "Foundation", "Neuromancer"}
	i := 0
	for {
		if i >= len(books) {
			break
		}
		fmt.Println(books[i])
		i++
	}

	// range-based — index and value together
	for i, v := range books {
		fmt.Println(i, v)
	}
}
