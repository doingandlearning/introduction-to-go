package main

import "fmt"

// concrete version — checking out a book that might not exist
func checkOut(catalog map[string]int, title string) (int, string, error) {

	if err := validateTitle(catalog, title); err != nil {
		return 0, "", err
	}

	if catalog[title] == 0 {
		return 0, "", fmt.Errorf("No copies of %s left in library", title)
	}
	catalog[title]--
	return catalog[title], title, nil
}

func validateTitle(catalog map[string]int, title string) error {
	if _, ok := catalog[title]; !ok {
		return fmt.Errorf("no such book: %s", title)
	}
	return nil
}

func main() {
	var catalog map[string]int
	remainingCopies, title, err := checkOut(catalog, "Dune")
	if err != nil {
		// i should handle ...

	}

}
