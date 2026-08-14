// Command connectors is Exercise 2: a simplified, single-method
// Abstract Factory for two storage backends, plus a function that
// accepts the factory interface generically.
package main

import "fmt"

// TODO (Exercise 2): declare a Store interface with one method,
// Describe() string.

// TODO (Exercise 2): declare a StoreFactory interface with one method,
// NewStore() Store.

// TODO (Exercise 2): implement two Store types - localStore and
// cloudStore - each with a Describe method returning a distinct string
// (e.g. "local disk store" / "cloud object store").

// TODO (Exercise 2): implement two factories - LocalFactory and
// CloudFactory - each satisfying StoreFactory by returning the matching
// Store type from NewStore().

// TODO (Exercise 2): implement describeBackend(f StoreFactory) that
// calls f.NewStore().Describe() and prints the result. It must not
// mention LocalFactory, CloudFactory, localStore, or cloudStore by name.

func main() {
	// TODO (Exercise 2): call describeBackend twice - once with
	// LocalFactory{}, once with CloudFactory{} - and confirm each prints
	// its own backend's description.
	fmt.Println("implement the TODOs above")
}
