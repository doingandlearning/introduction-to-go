// Exercise 5 (bonus): confirm sync.Once's guarantee by outcome, not by
// timing - after every goroutine has finished, every one of them must have
// received the exact same *Settings pointer.
package registry

import "testing"

func TestGetSettingsSingleInstanceUnderConcurrency(t *testing.T) {
	t.Skip("TODO: implement TestGetSettingsSingleInstanceUnderConcurrency")
}
