// Package domain holds the plain data types shared across every layer
// of the service. Nothing in this package knows about HTTP, storage, or
// business rules — it's just the shape of the data.
package domain

// Item is a single stock item tracked by the API.
type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
