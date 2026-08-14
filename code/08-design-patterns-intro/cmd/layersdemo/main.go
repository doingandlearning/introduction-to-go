// Command layersdemo builds the same Service twice: once wired to a
// "real" in-memory repository, once wired to a test-double repository
// that returns canned data. Service's code never changes - only what
// gets passed into NewService at the composition root.
//
//	go run ./cmd/layersdemo
package main

import (
	"fmt"

	"example.com/patterns-intro/internal/user"
)

// fakeRepository is a Repository implementation built purely for this
// demo, standing in for what a unit test's test double would look like.
type fakeRepository struct{}

func (fakeRepository) FindByID(id string) (*user.User, error) {
	return &user.User{ID: id, Name: "Test Double User"}, nil
}

func main() {
	// "Production" wiring: a real in-memory repository seeded with data.
	realRepo := user.NewInMemoryRepository(map[string]*user.User{
		"1": {ID: "1", Name: "Ada Lovelace"},
		"2": {ID: "2", Name: "Grace Hopper"},
	})
	prodService := user.NewService(realRepo)

	greeting, err := prodService.GetGreeting("1")
	if err != nil {
		fmt.Println("prod error:", err)
	} else {
		fmt.Println("prod service:", greeting)
	}

	// "Test" wiring: same Service constructor, a fake repository instead.
	// No framework, no mocking library - just a different value that
	// satisfies the same Repository interface.
	testService := user.NewService(fakeRepository{})

	greeting, err = testService.GetGreeting("anything")
	if err != nil {
		fmt.Println("test error:", err)
	} else {
		fmt.Println("test service:", greeting)
	}

	// Prove the failure path too: ask the real repo for a user that
	// doesn't exist.
	_, err = prodService.GetGreeting("does-not-exist")
	fmt.Println("expected not-found error:", err)
}
