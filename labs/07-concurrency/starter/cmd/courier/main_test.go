// Exercise 3: pre-written test for dispatchOrders. Run
// `go test ./cmd/courier/...` first — it fails until dispatchOrders is
// implemented below. If dispatchOrders never sends or never closes ch,
// this test can't fail with a normal assertion — nothing arrives, so it
// times out instead. That's expected and is itself a small lesson: a
// stuck goroutine doesn't report a failure on its own, only a hang.
package main

import (
	"testing"
	"time"
)

func TestDispatchOrders(t *testing.T) {
	ch := make(chan int)
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
