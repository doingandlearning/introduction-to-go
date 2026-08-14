// Command connectors is Exercise 2: a simplified, single-method
// Abstract Factory for two storage backends, plus a function that
// accepts the factory interface generically.
package main

import "fmt"

type Store interface {
	Describe() string
}

type StoreFactory interface {
	NewStore() Store
}

type localStore struct{}

func (localStore) Describe() string { return "local disk store" }

type cloudStore struct{}

func (cloudStore) Describe() string { return "cloud object store" }

type LocalFactory struct{}

func (LocalFactory) NewStore() Store { return localStore{} }

type CloudFactory struct{}

func (CloudFactory) NewStore() Store { return cloudStore{} }

// describeBackend is written entirely against StoreFactory - it never
// mentions a concrete factory or store type.
func describeBackend(f StoreFactory) {
	fmt.Println(f.NewStore().Describe())
}

func main() {
	describeBackend(LocalFactory{})
	describeBackend(CloudFactory{})
}
