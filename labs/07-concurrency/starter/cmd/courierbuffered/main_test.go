// Exercise 4: pre-written test for dispatchOrders. Run
// `go test ./cmd/courierbuffered/...` first — it fails until
// dispatchOrders is implemented below. This test only checks the values
// that arrive and that the channel closes — it does not (and can't
// reliably) assert anything about the "about to send" print timing;
// that's what task 3's `go run` comparison against Exercise 3 is for.
package main

import (
	"testing"
	"time"
)

func TestDispatchOrders(t *testing.T) {
	ch := make(chan int, 3)
	go dispatchOrders(ch)

	var got []int
loop:
	for {
		select {
		case orderID, ok := <-ch:
			if !ok {
				break loop
			}
			got = append(got, orderID)
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out waiting on ch — did dispatchOrders forget to send or close it? got so far: %v", got)
		}
	}

	want := []int{1, 2, 3, 4, 5}
	if len(got) != len(want) {
		t.Fatalf("got %d orders, want %d: %v", len(got), len(want), got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("order[%d] = %d, want %d", i, got[i], w)
		}
	}
}
