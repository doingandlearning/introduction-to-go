package spider

import "io"

// WriteCSV writes one CSV row per result to w: url, status, bytes,
// elapsed milliseconds, error (blank if none). The header row names
// every column — the general-purpose encoding/csv package, the same
// standard-library-first spirit as encoding/json back in Topic 10.
//
// TODO (Exercise 4): implement this. It should:
//  1. Add imports for "encoding/csv", "fmt", and "strconv" above
//     `package spider` — none of them are here yet.
//  2. csv.NewWriter(w), then Write a header row:
//     {"url", "status", "bytes", "elapsed_ms", "error"}.
//  3. For each Result, Write a row: r.URL, strconv.Itoa(r.StatusCode),
//     strconv.Itoa(r.Bytes),
//     strconv.FormatInt(r.Elapsed.Milliseconds(), 10), r.Err — wrap any
//     write error with fmt.Errorf.
//  4. cw.Flush(), then return cw.Error().
func WriteCSV(w io.Writer, results []Result) error {
	// TODO: replace this placeholder.
	return nil
}
