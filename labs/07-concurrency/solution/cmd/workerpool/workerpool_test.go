// Exercise 8: pre-written test for worker. Run
// `go test ./cmd/workerpool/...` first — it fails until worker (in
// main.go) is implemented.
//
// Testing concurrent code is different from testing a plain function:
// you cannot assert *which* worker processed job 7, or the exact order
// results arrive in results — the scheduler doesn't guarantee either, and
// a test that assumes one will flake sooner or later. What you *can*
// assert, every single run, is the aggregate: exactly 12 results come
// back, and the complete set of fees is {2, 4, 6, ..., 24} — that's true
// no matter which worker did the work or what order it finished in.
//
// If worker never touches jobs or never calls wg.Done(), this test can't
// fail with a normal assertion — wg.Wait() would just hang forever. The
// select below with a 2-second time.After guards against that: a stuck
// worker pool fails fast with a clear message instead of hanging until
// go test's own multi-minute default timeout kills the whole run.
package main

import (
	"sort"
	"sync"
	"testing"
	"time"
)

func TestWorkerPoolAggregateResults(t *testing.T) {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= 12; j++ {
		jobs <- j
	}
	close(jobs)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// every worker has drained jobs and exited
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for all workers to finish — did worker forget to range over jobs, or forget wg.Done()?")
	}

	close(results) // now safe — nobody is still sending to it

	var got []int
	for r := range results {
		got = append(got, r)
	}

	if len(got) != 12 {
		t.Fatalf("expected 12 results, got %d: %v", len(got), got)
	}

	sort.Ints(got) // arrival order isn't guaranteed, so compare sorted

	want := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("sorted result[%d] = %d, want %d (full sorted results: %v)", i, got[i], w, got)
		}
	}
}
