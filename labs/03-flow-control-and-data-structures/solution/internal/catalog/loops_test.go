// This test file is already complete — you're not writing it. It's the
// specification for Exercise 1: run `go test ./...` now, before touching
// loops.go, and all three tests fail. Implement PrintCatalog,
// CountAvailable, and FindFirstAvailable until they pass. Writing a test
// like this yourself is Topic 12's job, not this one's.
package catalog

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintCatalog(t *testing.T) {
	books := []Book{
		{Title: "The Go Programming Language", Copies: 2},
		{Title: "Design Patterns", Copies: 0},
	}

	var buf bytes.Buffer
	PrintCatalog(&buf, books)
	out := buf.String()

	for _, want := range []string{"The Go Programming Language", "Design Patterns"} {
		if !strings.Contains(out, want) {
			t.Errorf("PrintCatalog output is missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestCountAvailable(t *testing.T) {
	books := []Book{
		{Title: "A", Copies: 2},
		{Title: "B", Copies: 0},
		{Title: "C", Copies: 1},
		{Title: "D", Copies: 0},
	}

	got := CountAvailable(books)
	if got != 2 {
		t.Errorf("CountAvailable(...) = %d, want 2", got)
	}
}

func TestFindFirstAvailable(t *testing.T) {
	books := []Book{
		{Title: "A", Copies: 0},
		{Title: "B", Copies: 0},
		{Title: "C", Copies: 3},
	}

	got, ok := FindFirstAvailable(books)
	if !ok {
		t.Fatal("FindFirstAvailable(...) returned ok = false, want true")
	}
	if got.Title != "C" {
		t.Errorf("FindFirstAvailable(...) = %+v, want Title \"C\"", got)
	}

	none, ok := FindFirstAvailable([]Book{{Title: "X", Copies: 0}})
	if ok {
		t.Errorf("FindFirstAvailable(...) = %+v, ok = true, want ok = false", none)
	}
}
