package main

import "fmt"

type Point struct{ X, Y int }

func (p Point) String() string {
	return fmt.Sprintf("!!!(%d, %d)!!!", p.X, p.X)
} // __str__, __repr__

func main() {
	p := Point{X: 3, Y: 4}
	fmt.Println(p)
}
