package main

import "fmt"

func double(x *int) {
	*x = *x * 2
}

func main() {
	x := 5
	double(&x)
	fmt.Println(x)
}
