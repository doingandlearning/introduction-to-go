// Demo: nested JSON in a struct definition.
//
// Three shapes worth showing side by side:
//
//  1. A named nested struct field -> a genuinely nested JSON object.
//  2. An embedded (anonymous) struct -> its fields get PROMOTED and
//     flattened into the parent JSON object, not nested. This trips
//     people up who assume embedding always means nesting.
//  3. A slice of structs -> a JSON array of objects.
//
// Run:  go run ./cmd/nested-json
// Hit:  curl localhost:8080/item
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Dimensions is nested as a NAMED field below -> produces a real nested
// JSON object: {"dimensions": {"width_cm": ..., "height_cm": ...}}.
type Dimensions struct {
	WidthCM  float64 `json:"width_cm"`
	HeightCM float64 `json:"height_cm"`
}

// Audit is embedded ANONYMOUSLY below -> its fields are promoted and sit
// flat alongside Item's own fields in the JSON, NOT nested under "audit".
// This is the gotcha: embedding a struct type is not the same as nesting
// its JSON.
type Audit struct {
	CreatedBy string `json:"created_by"`
	Version   int    `json:"version"`
}

// Supplier is nested inside a slice below -> each element serializes as
// its own JSON object, all three sitting inside a JSON array.
type Supplier struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type Item struct {
	Name       string     `json:"name"`
	Dimensions Dimensions `json:"dimensions"` // named field -> nested object
	Audit                 // embedded, no field name -> fields promoted/flattened
	Suppliers  []Supplier `json:"suppliers"` // slice of structs -> JSON array of objects
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /item", func(w http.ResponseWriter, r *http.Request) {
		item := Item{
			Name: "widget",
			Dimensions: Dimensions{
				WidthCM:  12.5,
				HeightCM: 4,
			},
			Audit: Audit{
				CreatedBy: "kevin",
				Version:   1,
			},
			Suppliers: []Supplier{
				{Name: "Acme Ltd", Country: "UK"},
				{Name: "Global Parts", Country: "DE"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(item)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
