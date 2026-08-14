// This test file is already complete — you're not writing it. It's the
// specification for Exercises 2, 4, and 6: run `go test ./...` now,
// before touching main.go, and all three tests fail. Implement
// reportZeroValues, closeZonesFixed, and printItemFormats until they
// pass. Writing a test like this yourself is Topic 12's job, not this
// one's.
package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestReportZeroValues(t *testing.T) {
	var buf bytes.Buffer
	reportZeroValues(&buf)
	out := buf.String()

	for _, label := range []string{"int:", "float64:", "string:", "bool:", "*Item:", "Item:"} {
		if !strings.Contains(out, label) {
			t.Errorf("reportZeroValues output is missing a line for %q\ngot:\n%s", label, out)
		}
	}
}

func TestCloseZonesFixed(t *testing.T) {
	var buf bytes.Buffer
	closeZonesFixed(&buf)
	out := buf.String()

	// LIFO: the zone deferred last (C) should print first.
	wantOrder := []string{"closing zone: C", "closing zone: B", "closing zone: A"}
	lastIdx := -1
	for _, want := range wantOrder {
		idx := strings.Index(out, want)
		if idx == -1 {
			t.Fatalf("closeZonesFixed output is missing %q\ngot:\n%s", want, out)
		}
		if idx < lastIdx {
			t.Errorf("closeZonesFixed printed %q out of LIFO order\ngot:\n%s", want, out)
		}
		lastIdx = idx
	}
}

func TestPrintItemFormats(t *testing.T) {
	sample := Item{Name: "Widget", Quantity: 3, UnitPrice: 1.5}

	var buf bytes.Buffer
	printItemFormats(&buf, sample)
	out := buf.String()

	for _, want := range []string{"%v", "%+v", "%#v", "%T", "Widget", "main.Item"} {
		if !strings.Contains(out, want) {
			t.Errorf("printItemFormats output is missing %q\ngot:\n%s", want, out)
		}
	}
}
