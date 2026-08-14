// Exercise 8: prove the worker pool's aggregate result is deterministic,
// even though which worker processes which job is not.
package main

import (
	"sort"
	"sync"
	"testing"
)

// TestWorkerPoolAggregateResults sends a fixed set of jobs through the real
// worker function and checks the aggregate outcome: exactly 12 results,
// and the complete set of fees matches what job IDs 1-12 should produce.
// It never asserts an order — which worker handled which job, or which
// order results arrived in results, is not guaranteed by the scheduler
// and this test does not depend on it.
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

	wg.Wait()      // every worker has drained jobs and exited
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
