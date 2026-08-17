// Package urlfile reads a plain-text list of target URLs from disk —
// one URL per line, blank lines and lines starting with # ignored.
package urlfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadURLs opens path and returns every non-blank, non-comment line as
// a URL, in file order.
func ReadURLs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening url file: %w", err)
	}
	defer f.Close()

	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading url file: %w", err)
	}
	return urls, nil
}
