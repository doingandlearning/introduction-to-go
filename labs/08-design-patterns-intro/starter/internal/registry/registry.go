// Package registry holds process-wide application settings, loaded once
// and shared everywhere. This is Exercise 1-3's Singleton.
package registry

// Settings holds process-wide configuration.
type Settings struct {
	Environment string
}

var instance *Settings

// GetSettings returns the single, shared Settings instance.
//
// TODO (Exercise 1): implement the naive version of this function:
//   - If instance is nil, build one by calling loadSettings() and
//     assign it to instance.
//   - Return instance.
//
// This version will NOT be safe for concurrent use - that's the point
// of Exercise 2. You'll fix it with sync.Once in Exercise 3.
func GetSettings() *Settings {
	// TODO: replace this placeholder.
	return nil
}

// loadSettings simulates an expensive load (reading a file, hitting a
// secrets manager, etc). Do not modify this function.
func loadSettings() *Settings {
	return &Settings{Environment: "development"}
}
