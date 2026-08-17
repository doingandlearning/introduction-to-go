// Package urlfile reads a plain-text list of target URLs from disk —
// one URL per line, blank lines and lines starting with # ignored.
package urlfile

// ReadURLs opens path and returns every non-blank, non-comment line as
// a URL, in file order.
//
// TODO (Exercise 1): implement this. It should:
//  1. Add imports for "bufio", "fmt", "os", and "strings" above
//     `package urlfile` — none of them are here yet.
//  2. os.Open(path); wrap a failure with
//     fmt.Errorf("opening url file: %w", err).
//  3. defer f.Close().
//  4. bufio.NewScanner(f), Scan() in a loop, strings.TrimSpace each
//     line, and skip it if it's empty or starts with "#".
//  5. Check scanner.Err() after the loop and wrap it the same way as
//     step 2.
func ReadURLs(path string) ([]string, error) {
	// TODO: replace this placeholder.
	return nil, nil
}
