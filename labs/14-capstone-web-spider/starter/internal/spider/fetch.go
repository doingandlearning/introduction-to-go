// Package spider fetches a batch of URLs concurrently and reports what
// happened to each one.
package spider

import (
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
//
// TODO (Exercise 2): implement this. It should:
//  1. Add an `import "io"` line above `package spider` — this file
//     doesn't import it yet.
//  2. Record time.Now() before the request.
//  3. client.Get(url). On error, return a Result with URL, Err set to
//     err.Error(), and Elapsed already measured with time.Since.
//  4. Otherwise defer resp.Body.Close().
//  5. io.Copy(io.Discard, resp.Body) to read (and count) the whole
//     body — it returns the byte count and an error.
//  6. Return a Result with URL, resp.StatusCode, the byte count (as an
//     int), and Elapsed set — and Err set too, if the io.Copy itself
//     failed.
func Fetch(client *http.Client, url string) Result {
	// TODO: replace this placeholder.
	return Result{URL: url}
}
