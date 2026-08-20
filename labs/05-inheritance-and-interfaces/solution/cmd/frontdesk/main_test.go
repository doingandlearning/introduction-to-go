// This test file is already complete — you're not writing it. It's the
// specification for Exercise 4: run `go test ./...` now, before
// touching logCheckIn, and it fails — the placeholder writes the raw
// event with no type switch, so none of the expected substrings show
// up. Implement the type switch until it passes.
package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogCheckIn(t *testing.T) {
	var buf bytes.Buffer
	logCheckIn(&buf, 42)
	if got := buf.String(); !strings.Contains(got, "42") || !strings.Contains(got, "visitors") {
		t.Errorf("logCheckIn(42) output = %q, want it to mention 42 visitors", got)
	}

	buf.Reset()
	logCheckIn(&buf, "Priya")
	if got := buf.String(); !strings.Contains(got, "Priya") || !strings.Contains(got, "patron") {
		t.Errorf("logCheckIn(%q) output = %q, want it to mention patron Priya", "Priya", got)
	}

	buf.Reset()
	logCheckIn(&buf, 3.14)
	if got := buf.String(); !strings.Contains(got, "3.14") {
		t.Errorf("logCheckIn(3.14) output = %q, want it to mention the unrecognized value 3.14", got)
	}
}
