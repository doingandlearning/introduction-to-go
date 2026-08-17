// This test file is already complete — you're not writing it. It's the
// specification for Exercise 4: run `go test ./...` now, before touching
// report.go, and it fails. Implement WriteCSV until it passes. Writing
// a test like this yourself is Topic 12's job for every exercise in
// this lab except the last one.
package spider

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteCSV(t *testing.T) {
	results := []Result{
		{URL: "https://example.com", StatusCode: 200, Bytes: 512, Elapsed: 120 * time.Millisecond},
		{URL: "https://broken.example", Err: "connection refused"},
	}

	var buf bytes.Buffer
	if err := WriteCSV(&buf, results); err != nil {
		t.Fatalf("WriteCSV(...) unexpected error: %v", err)
	}
	out := buf.String()

	for _, want := range []string{
		"url", "status", "bytes", "elapsed_ms", "error",
		"https://example.com", "200", "512",
		"https://broken.example", "connection refused",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("WriteCSV output is missing %q\ngot:\n%s", want, out)
		}
	}

	if strings.Count(out, "\n") < 3 {
		t.Errorf("WriteCSV output has fewer than 3 lines (header + 2 rows)\ngot:\n%s", out)
	}
}
