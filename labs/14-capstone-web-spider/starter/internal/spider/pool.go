package spider

// Run fetches every url in urls concurrently, using a fixed pool of
// workers, and returns one Result per url — the same worker-pool shape
// as Topic 7's lecture example, generalized: fetch is a parameter
// instead of a hardcoded call, so this can be tested without a real
// network or even a real HTTP client.
//
// TODO (Exercise 3): implement this. It should:
//  1. Add an `import "sync"` line above `package spider` — this file
//     doesn't import it yet.
//  2. Make a buffered `jobs chan string` and `results chan Result`,
//     both sized len(urls).
//  3. Start `workers` goroutines, each looping `for url := range jobs`,
//     calling fetch(url), and sending the Result to results — wrapped
//     in a sync.WaitGroup so you know when every worker has exited.
//  4. Send every url into jobs, then close(jobs).
//  5. wg.Wait(), then close(results).
//  6. Drain results into a []Result and return it.
func Run(urls []string, workers int, fetch func(url string) Result) []Result {
	// TODO: replace this placeholder.
	return nil
}
