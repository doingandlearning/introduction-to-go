// This test file is already complete — you're not writing it. It's the
// specification for Exercise 3: run `go test ./...` now, before touching
// pool.go, and it fails. Implement Run until it passes. This test
// asserts the AGGREGATE result, never which worker handled which url or
// what order results came back in — exactly the principle Topic 7's own
// summary named: "test outcomes, not order." Writing a test like this
// yourself is Topic 12's job for every exercise in this lab except the
// last one.
package spider

import (
	"sort"
	"testing"
)

func TestRun(t *testing.T) {
	urls := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}

	fake := func(url string) Result {
		return Result{URL: url, StatusCode: 200}
	}

	got := Run(urls, 3, fake)

	if len(got) != len(urls) {
		t.Fatalf("Run(...) returned %d results, want %d", len(got), len(urls))
	}

	var gotURLs []string
	for _, r := range got {
		gotURLs = append(gotURLs, r.URL)
	}
	sort.Strings(gotURLs)

	wantURLs := append([]string(nil), urls...)
	sort.Strings(wantURLs)

	for i := range wantURLs {
		if gotURLs[i] != wantURLs[i] {
			t.Errorf("Run(...) result urls = %v, want (in some order) %v", gotURLs, wantURLs)
			break
		}
	}
}

func TestRun_EmptyInput(t *testing.T) {
	fake := func(url string) Result { return Result{URL: url} }

	got := Run(nil, 3, fake)

	if len(got) != 0 {
		t.Errorf("Run(nil, ...) = %v, want empty", got)
	}
}
