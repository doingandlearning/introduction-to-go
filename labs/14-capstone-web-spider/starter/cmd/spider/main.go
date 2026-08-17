// Command spider fetches a list of URLs concurrently and writes a CSV
// report of what happened to each one.
//
//	go run ./cmd/spider -urls sample-urls.txt -out report.csv -workers 4
package main

// TODO (Exercise 4): implement main(). It should:
//  1. Add imports for "flag", "fmt", "net/http", "os", "time",
//     "example.com/spider/internal/spider", and
//     "example.com/spider/internal/urlfile" above `package main` —
//     none of them are here yet.
//  2. flag.String("urls", "sample-urls.txt", ...),
//     flag.String("out", "report.csv", ...),
//     flag.Int("workers", 4, ...), and
//     flag.Duration("timeout", 5*time.Second, ...) — then
//     flag.Parse().
//  3. urlfile.ReadURLs(*urlsPath); exit with a message on error or on
//     an empty list (fmt.Fprintln(os.Stderr, ...); os.Exit(1)).
//  4. Build a &http.Client{Timeout: *timeout}, and call spider.Run with
//     a closure that calls spider.Fetch(client, url).
//  5. os.Create(*outPath); defer its Close(); pass it to
//     spider.WriteCSV(out, results).
//  6. Print a one-line summary of what happened.
func main() {
	// TODO: replace this placeholder — see the steps above.
	panic("TODO: implement main")
}
