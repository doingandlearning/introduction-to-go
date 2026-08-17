# Lab 14 (Capstone): a concurrent URL-status spider

Starter code is in `starter/` (TODOs to fill in). A complete reference is
in `solution/` — don't look until you've had a go.

This lab doesn't teach one new language feature — it combines four
things into one small, realistic tool: real file I/O (new today), a
real HTTP client (new today), CLI flags (new today), and the worker-pool
concurrency pattern from **Topic 7**. The scenario: a small internal
tool that reads a list of URLs from a file, checks every one
concurrently, and writes a CSV report of what happened to each.

**Every exercise except the last ships with its test already written**,
the same pattern as every lab since Topic 2 — sitting in
`starter/internal/urlfile/urlfile_test.go`,
`starter/internal/spider/fetch_test.go`, `pool_test.go`, and
`report_test.go`. Run `go test ./...` from `starter/` right now — all
four fail. **Exercise 5 is different: you write the test yourself**, for
the first time since Topic 12 unlocked that skill — the one place in
this lab where the old "write your own test" pattern comes back, on
purpose.

None of the tests in this lab touch the real internet — every HTTP test
uses `httptest.NewServer`, the same tool Topic 10 used to test a
handler, so they run identically whether you're online or not.
`sample-urls.txt` (real internet addresses, including one that's
deliberately unreachable) is only for the `go run` verification steps,
not for grading.

Both directories share the same shape:

```
cmd/spider/main.go                     composition root: flags, wiring, the CSV file
internal/urlfile/urlfile.go            Exercise 1: read a URL list from a file
internal/spider/fetch.go               Exercise 2: fetch one URL
internal/spider/pool.go                Exercise 3: fetch many URLs concurrently
internal/spider/report.go              Exercise 4: write a CSV report
internal/spider/integration_test.go    Exercise 5: your test
```

---

## Exercise 1: Read the URL list — real file I/O

**Objective:** Open and scan a real file on disk — something every
earlier lab's `cmd/*` package implied was possible but never actually
did.

**Context:** `starter/internal/urlfile/urlfile.go` has a `TODO` for
`ReadURLs`. `TestReadURLs` and `TestReadURLs_MissingFile` in
`urlfile_test.go` are already written and already failing.

**Tasks:**

1. Run `go test ./...` from `starter/`. Read both failures — one checks
   a well-formed file (skipping blank lines and `#` comments), the
   other checks that a missing file returns a real error instead of a
   panic or a silent empty result.
2. Implement `ReadURLs(path string) ([]string, error)` in
   `urlfile.go`. The file doesn't import `"bufio"`, `"fmt"`, `"os"`, or
   `"strings"` yet — add all four yourself. `os.Open`, wrap a failure
   with `fmt.Errorf("...: %w", err)`, `defer f.Close()`,
   `bufio.NewScanner(f)`, and skip blank/`#`-prefixed lines after
   `strings.TrimSpace`.
3. Re-run `go test ./...` and confirm both tests pass.

**Key Learning:** `os.Open` plus `bufio.Scanner` is the default way to
read a text file line by line in Go — no special "read a file" library,
just the same `io.Reader` machinery you've already seen backing HTTP
request bodies and test buffers throughout the course.

---

## Exercise 2: Fetch one URL — a real HTTP client

**Objective:** Make an outbound HTTP request for the first time in this
course — every earlier topic had you build the *server* side only.

**Context:** `starter/internal/spider/fetch.go` has a `TODO` for
`Fetch`. `TestFetch_Success`, `TestFetch_NonOKStatus`, and
`TestFetch_ConnectionError` in `fetch_test.go` are already written and
already failing — all three run against `httptest.NewServer`, the exact
tool Topic 10 used to test a handler without a real network call; here
it stands in for the *remote* server instead.

**Tasks:**

1. Run `go test ./...`. Read all three failures — one checks a normal
   200 response, one checks that a 404 is still a successful `Result`
   (not a returned `error`), one closes the test server first to force
   a connection failure deterministically.
2. Implement `Fetch(client *http.Client, url string) Result` in
   `fetch.go`. The file doesn't import `"io"` yet — add that line
   yourself. Time the call with `time.Now()`/`time.Since`, call
   `client.Get(url)`, and on error return a `Result` with `Err` set —
   never a returned `error`; a failed fetch is still data, not a
   program failure. Otherwise `defer resp.Body.Close()`,
   `io.Copy(io.Discard, resp.Body)` to read and count the body, and
   return a populated `Result`.
3. Re-run `go test ./...` and confirm all three pass.

**Key Learning:** `*http.Client` is the client-side mirror of everything
you already know about `net/http` from the server side. A non-2xx
status is not a Go error — checking `resp.StatusCode` yourself is the
same "no automatic exception translation" philosophy Topic 10 named for
handlers, applying just as much to the client making the call.

---

## Exercise 3: Fetch many URLs concurrently — the worker pool comes back

**Objective:** Reuse Topic 7's worker pool shape, generalized so it
never has to know what "fetch" actually means.

**Context:** `starter/internal/spider/pool.go` has a `TODO` for `Run`.
`TestRun` and `TestRun_EmptyInput` in `pool_test.go` are already written
and already failing. `Run` takes `fetch` as a parameter — a function
value, the tool from **Topic 6** — so its test passes in a fake `fetch`
and needs no real HTTP client or network at all.

**Tasks:**

1. Run `go test ./...`. Read `TestRun`'s failure — notice it sorts both
   sides before comparing and never checks which worker handled which
   URL or what order results came back in. That's deliberate: **"test
   outcomes, not order"** is Topic 7's own summary line, and this is
   where you apply it instead of just reading it.
2. Implement `Run(urls []string, workers int, fetch func(url string)
   Result) []Result` in `pool.go`. The file doesn't import `"sync"`
   yet — add that line yourself. Same shape as Topic 7's lecture
   example: buffered `jobs`/`results` channels, a fixed crew of
   goroutines ranging over `jobs`, a `sync.WaitGroup` gating when it's
   safe to close `results`.
3. Re-run `go test ./...` and confirm both tests pass.

**Key Learning:** A worker pool doesn't need to know what kind of work
it's doing. Passing `fetch` in as a function value is what makes this
one testable with zero network calls, and reusable for work that isn't
HTTP at all.

---

## Exercise 4: Write the CSV report and wire up the CLI

**Objective:** Write a second, more general-purpose file format — this
time one your program produces — and parse real command-line flags for
the first time in the course.

**Context:** `starter/internal/spider/report.go` has a `TODO` for
`WriteCSV`. `TestWriteCSV` in `report_test.go` is already written and
already failing. `starter/cmd/spider/main.go` has a `TODO` for the
whole `main()` — the composition root, wiring flags, `urlfile.ReadURLs`,
`spider.Run`, and `spider.WriteCSV` together, the same job Topics 8-10's
`main` packages did.

**Tasks:**

1. Run `go test ./...`. Read `TestWriteCSV`'s failure — it checks for
   the header row and both a successful and a failed result's fields as
   substrings, not one exact byte-for-byte layout.
2. Implement `WriteCSV(w io.Writer, results []Result) error` in
   `report.go`, using `encoding/csv` — add `"encoding/csv"`, `"fmt"`,
   and `"strconv"` yourself, they're not imported yet.
3. Re-run `go test ./...` and confirm `TestWriteCSV` passes.
4. Implement `main()` in `cmd/spider/main.go`: `flag.String`/`flag.Int`/
   `flag.Duration` for the URL file path, the output path, worker
   count, and per-request timeout; `flag.Parse()`; then wire
   `urlfile.ReadURLs` -> `spider.Run` -> `os.Create` -> `spider.WriteCSV`
   together, exiting with a clear message on any error.
5. Run `go run ./cmd/spider -urls sample-urls.txt -out report.csv
   -workers 4` from `starter/` and open `report.csv` — confirm it has a
   row for every URL in `sample-urls.txt`, including the one that's
   deliberately unreachable.

**Key Learning:** `flag` is the standard library's answer to CLI
arguments — no separate package to choose, the same "there is no
decision" theme Topic 12 named for testing frameworks. And
`encoding/csv` proves the streaming-vs-general-purpose split from Topic
10's JSON slide isn't JSON-specific — it's just how Go's
`io.Writer`-based encoders work.

---

## Exercise 5: Write your own test

**Objective:** For the first time since Topic 12, write a test from a
blank function — no pre-written specification this time, just a job to
do and the tools from every earlier exercise in this lab to do it with.

**Context:** `starter/internal/spider/integration_test.go` has
`TestEndToEnd` as a `t.Skip` stub — the same shape Topic 12's lab used,
not the pre-written-and-failing pattern every other exercise in this lab
(and every lab before Topic 12) used.

**Tasks:**

1. Delete the `t.Skip(...)` line.
2. Spin up an `httptest.NewServer` (like `fetch_test.go` does) that
   returns a fixed, recognizable response.
3. Write a real URL file into a `t.TempDir()` (like `urlfile_test.go`
   does) pointing at that test server, using `os.WriteFile`.
4. Call `urlfile.ReadURLs`, then `spider.Run` with `spider.Fetch` wired
   to the test server's client, then `spider.WriteCSV` to a real file
   in the same temp directory.
5. Read the report file back with `os.ReadFile` and assert it contains
   what you'd expect from the fake server's response.
6. Run `go test ./...` and confirm `TestEndToEnd` passes alongside every
   other test in this lab.

**Key Learning:** Every tool this test needs — `httptest.NewServer`,
`t.TempDir`, a multi-step assertion with no framework — was already
introduced somewhere earlier in the course. Writing this one from
scratch is the actual payoff of Topic 12, not a new topic on top of it.

---

## Summary

By the end of this lab you should be able to:

- Read a text file line by line with `os.Open` and `bufio.Scanner`
- Make an outbound HTTP request with `*http.Client`, and treat a
  non-2xx status as data, not an error
- Reuse a worker pool for work that has nothing to do with the original
  example, by passing the work itself in as a function value
- Write a second file format (`encoding/csv`) using the same
  `io.Writer` shape as everything else in the course
- Parse real command-line flags with the standard library's `flag`
  package
- Write a test completely from scratch — the actual skill Topic 12
  unlocked, applied here for the first time since that topic ended
