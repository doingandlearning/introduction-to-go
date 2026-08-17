// Command spider fetches a list of URLs concurrently and writes a CSV
// report of what happened to each one.
//
//	go run ./cmd/spider -urls sample-urls.txt -out report.csv -workers 4
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/spider/internal/spider"
	"example.com/spider/internal/urlfile"
)

func main() {
	urlsPath := flag.String("urls", "sample-urls.txt", "path to a file of URLs, one per line")
	outPath := flag.String("out", "report.csv", "path to write the CSV report to")
	workers := flag.Int("workers", 4, "number of concurrent workers")
	timeout := flag.Duration("timeout", 5*time.Second, "per-request timeout")
	flag.Parse()

	urls, err := urlfile.ReadURLs(*urlsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading url file:", err)
		os.Exit(1)
	}
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "no urls found in", *urlsPath)
		os.Exit(1)
	}

	client := &http.Client{Timeout: *timeout}
	results := spider.Run(urls, *workers, func(url string) spider.Result {
		return spider.Fetch(client, url)
	})

	out, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "creating report file:", err)
		os.Exit(1)
	}
	defer out.Close()

	if err := spider.WriteCSV(out, results); err != nil {
		fmt.Fprintln(os.Stderr, "writing report:", err)
		os.Exit(1)
	}

	fmt.Printf("fetched %d urls with %d workers -> %s\n", len(urls), *workers, *outPath)
}
