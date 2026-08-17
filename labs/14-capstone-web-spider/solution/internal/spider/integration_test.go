package spider

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"example.com/spider/internal/urlfile"
)

// TestEndToEnd wires urlfile.ReadURLs -> Run -> WriteCSV together
// against a real httptest.Server and a real temp file — the same
// integration a student writes for Exercise 5, now as the reference
// answer.
func TestEndToEnd(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	urlsPath := filepath.Join(dir, "urls.txt")
	content := srv.URL + "\n" + srv.URL + "/other\n"
	if err := os.WriteFile(urlsPath, []byte(content), 0o644); err != nil {
		t.Fatalf("writing fixture url file: %v", err)
	}

	urls, err := urlfile.ReadURLs(urlsPath)
	if err != nil {
		t.Fatalf("ReadURLs(%q): %v", urlsPath, err)
	}

	client := srv.Client()
	results := Run(urls, 2, func(url string) Result {
		return Fetch(client, url)
	})

	reportPath := filepath.Join(dir, "report.csv")
	out, err := os.Create(reportPath)
	if err != nil {
		t.Fatalf("creating report file: %v", err)
	}
	if err := WriteCSV(out, results); err != nil {
		out.Close()
		t.Fatalf("WriteCSV(...): %v", err)
	}
	out.Close()

	got, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("reading report file back: %v", err)
	}

	for _, want := range []string{srv.URL, "200"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("report file is missing %q\ngot:\n%s", want, got)
		}
	}
	if strings.Count(string(got), "200") < 2 {
		t.Errorf("expected both fetches to report status 200, got:\n%s", got)
	}
}
