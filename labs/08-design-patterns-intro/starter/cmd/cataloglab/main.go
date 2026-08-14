// Command cataloglab exercises the catalog Service with two different
// Repository implementations. See labs/exercise.md, Exercise 4.
//
//	go run ./cmd/cataloglab
package main

import (
	"fmt"

	"example.com/patterns-lab/internal/catalog"
)

// fakeRepository is a hand-written test double - no framework involved.
type fakeRepository struct{}

func (fakeRepository) FindByISBN(isbn string) (*catalog.Book, error) {
	return &catalog.Book{ISBN: isbn, Title: "Fake Book (test double)"}, nil
}

func main() {
	// "Production" wiring.
	realRepo := catalog.NewInMemoryRepository(map[string]*catalog.Book{
		"978-0134190440": {ISBN: "978-0134190440", Title: "The Go Programming Language"},
	})
	prodService := catalog.NewService(realRepo)

	desc, err := prodService.Describe("978-0134190440")
	if err != nil {
		fmt.Println("prod error:", err)
	} else {
		fmt.Println("prod:", desc)
	}

	// "Test" wiring - same Service constructor, a fake Repository.
	testService := catalog.NewService(fakeRepository{})

	desc, err = testService.Describe("anything")
	if err != nil {
		fmt.Println("test error:", err)
	} else {
		fmt.Println("test:", desc)
	}
}
