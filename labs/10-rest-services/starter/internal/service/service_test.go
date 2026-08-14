package service

import "testing"

// TestItemService_Create_StoresValidItem should construct a real
// repository.NewInMemoryRepository(), wrap it in NewItemService, call
// Create with a valid item, and confirm the item comes back with a
// generated ID and lands in the repository.
func TestItemService_Create_StoresValidItem(t *testing.T) {
	t.Skip("TODO: implement TestItemService_Create_StoresValidItem")
}

// TestItemService_Create_RejectsInvalidQuantity should call Create with
// Quantity: 0 and confirm the returned error satisfies
// errors.Is(err, ErrValidation).
func TestItemService_Create_RejectsInvalidQuantity(t *testing.T) {
	t.Skip("TODO: implement TestItemService_Create_RejectsInvalidQuantity")
}
