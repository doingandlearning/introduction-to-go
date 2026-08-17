package spider

import "sync"

// Run fetches every url in urls concurrently, using a fixed pool of
// workers, and returns one Result per url — the same worker-pool shape
// as Topic 7's lecture example, generalized: fetch is a parameter
// instead of a hardcoded call, so this can be tested without a real
// network or even a real HTTP client.
func Run(urls []string, workers int, fetch func(url string) Result) []Result {
	jobs := make(chan string, len(urls))
	results := make(chan Result, len(urls))
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				results <- fetch(url)
			}
		}()
	}

	for _, u := range urls {
		jobs <- u
	}
	close(jobs)

	wg.Wait()
	close(results)

	out := make([]Result, 0, len(urls))
	for r := range results {
		out = append(out, r)
	}
	return out
}
