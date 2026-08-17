// This test file is already complete — you're not writing it. It's the
// specification for Exercise 1: run `go test ./...` now, before touching
// urlfile.go, and both tests fail. Implement ReadURLs until they pass.
// Writing a test like this yourself is Topic 12's job for every
// exercise in this lab except the last one.
package urlfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadURLs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "urls.txt")
	content := "https://example.com\n\n# a comment, ignored\nhttps://go.dev\n   \nhttps://httpbin.org/status/200\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing fixture file: %v", err)
	}

	got, err := ReadURLs(path)
	if err != nil {
		t.Fatalf("ReadURLs(%q) unexpected error: %v", path, err)
	}

	want := []string{"https://example.com", "https://go.dev", "https://httpbin.org/status/200"}
	if len(got) != len(want) {
		t.Fatalf("ReadURLs(%q) = %v, want %v", path, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ReadURLs(%q)[%d] = %q, want %q", path, i, got[i], want[i])
		}
	}
}

func TestReadURLs_MissingFile(t *testing.T) {
	_, err := ReadURLs(filepath.Join(t.TempDir(), "does-not-exist.txt"))
	if err == nil {
		t.Fatal("ReadURLs on a missing file: expected an error, got nil")
	}
}
