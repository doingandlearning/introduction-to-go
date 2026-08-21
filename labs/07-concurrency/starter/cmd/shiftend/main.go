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
func checkOutCourier(w io.Writer, courierID int, wg *sync.WaitGroup) {
	// TODO(Exercise 2):
	//   1. defer wg.Done() first, same pairing as Exercise 1.
	//   2. Fprintf(w, "courier <id>: radio open"), then immediately defer
	//      an Fprintf of "courier <id>: radio closed".
	//   3. Fprintf(w, "courier <id>: route log started"), then
	//      immediately defer an Fprintf of "courier <id>: route log
	//      filed".
	//   4. If courierID == 3, Fprintf(w, "courier 3: route blocked,
	//      aborting") and return right here — before the delivery print
	//      below. Both deferred prints must still run.
	//   5. Fprintf(w, "courier <id>: delivering final parcel").
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
