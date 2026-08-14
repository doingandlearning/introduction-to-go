// Package config demonstrates the Singleton pattern in Go, first the
// naive (unsafe) version, then the idiomatic sync.Once fix.
package config

import (
	"os"
	"sync"
)

// Config holds process-wide configuration. In real code this might be
// populated from environment variables, a config file, or a secrets
// manager - the loading itself is what we want to guarantee happens
// exactly once.
type Config struct {
	APIKey string
}

var (
	instance    *Config
	once        sync.Once
	loadCount   int
	loadCountMu sync.Mutex
)

// GetConfig returns the single, shared Config instance. The first call
// builds it; every subsequent call - from any goroutine - gets back the
// same pointer.
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{APIKey: loadFromEnv()}
	})
	return instance
}

// LoadCount reports how many times loadFromEnv has actually run. Used by
// the demo and the lab to prove sync.Once only ever runs the loader once,
// no matter how many goroutines call GetConfig concurrently.
func LoadCount() int {
	loadCountMu.Lock()
	defer loadCountMu.Unlock()
	return loadCount
}

// loadFromEnv simulates an expensive or side-effecting load - reading an
// environment variable, hitting a secrets manager, whatever it might be
// in a real service.
func loadFromEnv() string {
	loadCountMu.Lock()
	loadCount++
	loadCountMu.Unlock()

	key := os.Getenv("DEMO_API_KEY")
	if key == "" {
		key = "dev-local-key"
	}
	return key
}
