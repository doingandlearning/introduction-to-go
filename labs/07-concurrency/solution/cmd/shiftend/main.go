// Exercise 2: use defer to guarantee each courier's radio and route log get
// closed out, however checkOutCourier returns.
package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// checkOutCourier takes an io.Writer instead of printing straight to
// stdout so the pre-written tests in main_test.go can capture its output
// with a bytes.Buffer and assert on it. main passes os.Stdout below.
//
// With 5 goroutines running checkOutCourier at once, deleting both defer
// statements and moving the "closed"/"filed" prints to the end of the
// function would mean any goroutine that takes the early-return path (here,
// courier 3) skips them entirely — its radio and route log never get
// closed out. That's 5 independent chances to leak a resource, not one:
// defer ties the cleanup to entry into the function, not to reaching a
// particular line at the end of it.
func checkOutCourier(w io.Writer, courierID int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Fprintf(w, "courier %d: radio open\n", courierID)
	defer fmt.Fprintf(w, "courier %d: radio closed\n", courierID)

	fmt.Fprintf(w, "courier %d: route log started\n", courierID)
	defer fmt.Fprintf(w, "courier %d: route log filed\n", courierID)

	if courierID == 3 {
		fmt.Fprintln(w, "courier 3: route blocked, aborting")
		return
	}

	fmt.Fprintf(w, "courier %d: delivering final parcel\n", courierID)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go checkOutCourier(os.Stdout, i, &wg)
	}

	wg.Wait()
	fmt.Println("all couriers checked out")
}
