package main

import "fmt"

func init() {
	fmt.Println("Starting pacakge ....")
}

func init() {
	fmt.Println("I'm really starting now!")
}

func main() {

	type Item struct {
		Name  string
		count int
	}

	i := Item{Name: "Widget", count: 3}
	fmt.Println(i.Name, i.count)
}

type Item struct {
	Name  string
	count int
}

func NewItem(name string) Item {
	i := Item{Name: name, count: 1}
	return i
}
