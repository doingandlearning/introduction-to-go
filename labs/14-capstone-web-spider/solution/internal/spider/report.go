package spider

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// WriteCSV writes one CSV row per result to w: url, status, bytes,
// elapsed milliseconds, error (blank if none). The header row names
// every column — the general-purpose encoding/csv package, the same
// standard-library-first spirit as encoding/json back in Topic 10.
func WriteCSV(w io.Writer, results []Result) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"url", "status", "bytes", "elapsed_ms", "error"}); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}
	for _, r := range results {
		row := []string{
			r.URL,
			strconv.Itoa(r.StatusCode),
			strconv.Itoa(r.Bytes),
			strconv.FormatInt(r.Elapsed.Milliseconds(), 10),
			r.Err,
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("writing row for %s: %w", r.URL, err)
		}
	}
	cw.Flush()
	return cw.Error()
}
