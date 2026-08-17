// Package spider fetches a batch of URLs concurrently and reports what
// happened to each one.
package spider

import (
	"io"
	"net/http"
	"time"
)

// Result is what happened when Fetch tried one URL. Err is a plain
// string, not an error — a fetch that fails still needs to end up in
// the results slice, not stop the batch or vanish.
type Result struct {
	URL        string
	StatusCode int
	Bytes      int
	Elapsed    time.Duration
	Err        string
}

// Fetch requests url with client and reports what happened. A network
// error or a non-2xx status is not itself fatal — it's just data: it
// comes back as a Result with Err set (or a StatusCode outside 2xx),
// never as a returned error.
func Fetch(client *http.Client, url string) Result {
	start := time.Now()

	resp, err := client.Get(url)
	if err != nil {
		return Result{URL: url, Err: err.Error(), Elapsed: time.Since(start)}
	}
	defer resp.Body.Close()

	n, err := io.Copy(io.Discard, resp.Body)
	elapsed := time.Since(start)
	if err != nil {
		return Result{URL: url, StatusCode: resp.StatusCode, Err: err.Error(), Elapsed: elapsed}
	}

	return Result{URL: url, StatusCode: resp.StatusCode, Bytes: int(n), Elapsed: elapsed}
}
