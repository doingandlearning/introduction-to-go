// Package registry holds process-wide application settings, loaded once
// and shared everywhere. This is Exercise 1-3's Singleton, fixed with
// sync.Once.
package registry

import "sync"

// Settings holds process-wide configuration.
type Settings struct {
	Environment string
}

var (
	instance *Settings
	once     sync.Once
)

// GetSettings returns the single, shared Settings instance. The first
// call builds it; every subsequent call - from any goroutine - gets back
// the same pointer, and loadSettings runs exactly once.
func GetSettings() *Settings {
	once.Do(func() {
		instance = loadSettings()
	})
	return instance
}

// loadSettings simulates an expensive load (reading a file, hitting a
// secrets manager, etc).
func loadSettings() *Settings {
	return &Settings{Environment: "development"}
}
