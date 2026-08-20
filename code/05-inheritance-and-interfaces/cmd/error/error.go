package main

import "fmt"

type MyError struct{}

func (e *MyError) Error() string { return "something broke" }

func doWork() error {
	if 1 > 2 {
		return &MyError{}
	}
	return nil
}

// func doWork() error {
// 	var e *MyError = nil
// 	// ... some logic that never sets e ...
// 	return e
// }

func main() {
	err := doWork()
	fmt.Printf("%T", err)
	fmt.Printf("%T", nil)
	if err == nil {
		fmt.Println("all good")
	} else {
		fmt.Println("failed:", err)
	}
}
