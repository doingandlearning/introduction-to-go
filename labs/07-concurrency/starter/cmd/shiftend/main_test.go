// Exercise 2: pre-written tests for checkOutCourier. Run
// `go test ./cmd/shiftend/...` first — it fails until checkOutCourier is
// implemented below. These tests prove that defer's cleanup still runs,
// in LIFO order, even on the early-return path.
package main

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

// TestCheckOutCourierEarlyReturn proves both deferred cleanup prints
// still run for courier 3, even though it returns before the delivery
// line.
func TestCheckOutCourierEarlyReturn(t *testing.T) {
	var buf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	checkOutCourier(&buf, 3, &wg)

	out := buf.String()

	if !strings.Contains(out, "courier 3: route blocked, aborting") {
		t.Errorf("expected the early-return message, got:\n%s", out)
	}
	if !strings.Contains(out, "courier 3: radio closed") {
		t.Errorf("expected the radio to still close on early return, got:\n%s", out)
	}
	if !strings.Contains(out, "courier 3: route log filed") {
		t.Errorf("expected the route log to still be filed on early return, got:\n%s", out)
	}
	if strings.Contains(out, "delivering final parcel") {
		t.Errorf("courier 3 should never reach the delivery line, got:\n%s", out)
	}
}

// TestCheckOutCourierFullRun proves a non-blocked courier's five lines
// come out in the right order — and specifically that defer's LIFO order
// puts "route log filed" before "radio closed", even though the radio
// opened first.
func TestCheckOutCourierFullRun(t *testing.T) {
	var buf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	checkOutCourier(&buf, 1, &wg)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	want := []string{
		"courier 1: radio open",
		"courier 1: route log started",
		"courier 1: delivering final parcel",
		"courier 1: route log filed",
		"courier 1: radio closed",
	}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d:\n%s", len(lines), len(want), buf.String())
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d = %q, want %q", i, lines[i], w)
		}
	}
}
