// This test file is already complete — you're not writing it. It's the
// specification for Exercise 2: run `go test ./...` now, before touching
// fetch.go, and it fails. Implement Fetch until it passes. All three
// cases run against httptest.NewServer — the exact tool Topic 10 used
// to test a handler without a real network call; here it stands in for
// the remote server instead, so nothing here needs real internet
// access. Writing a test like this yourself is Topic 12's job for every
// exercise in this lab except the last one.
package spider

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch_Success(t *testing.T) {
	body := "hello from the test server"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer srv.Close()

	got := Fetch(srv.Client(), srv.URL)

	if got.Err != "" {
		t.Fatalf("Fetch(%q) unexpected Err: %q", srv.URL, got.Err)
	}
	if got.StatusCode != http.StatusOK {
		t.Errorf("Fetch(%q).StatusCode = %d, want %d", srv.URL, got.StatusCode, http.StatusOK)
	}
	if got.Bytes != len(body) {
		t.Errorf("Fetch(%q).Bytes = %d, want %d", srv.URL, got.Bytes, len(body))
	}
}

func TestFetch_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	got := Fetch(srv.Client(), srv.URL)

	if got.Err != "" {
		t.Fatalf("Fetch(%q) unexpected Err: %q — a non-2xx status is still a successful fetch, not an error", srv.URL, got.Err)
	}
	if got.StatusCode != http.StatusNotFound {
		t.Errorf("Fetch(%q).StatusCode = %d, want %d", srv.URL, got.StatusCode, http.StatusNotFound)
	}
}

// TestFetch_ConnectionError closes the test server before Fetch ever
// reaches it, so the connection failure is deterministic and needs no
// real network access.
func TestFetch_ConnectionError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close()

	got := Fetch(http.DefaultClient, url)

	if got.Err == "" {
		t.Fatal("Fetch against a closed server: expected Err to be set, got empty string")
	}
	if got.URL != url {
		t.Errorf("Fetch(...).URL = %q, want %q", got.URL, url)
	}
}
