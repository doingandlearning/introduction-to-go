// Exercise 5 (bonus): confirm sync.Once's guarantee by outcome, not by
// timing - after every goroutine has finished, every one of them must have
// received the exact same *Settings pointer.
package registry

import (
	"sync"
	"testing"
)

func TestGetSettingsSingleInstanceUnderConcurrency(t *testing.T) {
	const goroutines = 50
	pointers := make([]*Settings, goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			pointers[i] = GetSettings()
		}()
	}
	wg.Wait()

	first := pointers[0]
	if first == nil {
		t.Fatal("GetSettings() returned nil")
	}
	for i, p := range pointers {
		if p != first {
			t.Errorf("goroutine %d got pointer %p, want %p - sync.Once should guarantee a single instance", i, p, first)
		}
	}
}
