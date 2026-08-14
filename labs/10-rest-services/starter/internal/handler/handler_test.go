package handler

import "testing"

// TestItemHandler_Create_ReturnsCreated should build a real service
// (real service.NewItemService over a real repository.NewInMemoryRepository),
// wrap it in NewItemHandler, send a request built with httptest.NewRequest,
// record the response with httptest.NewRecorder, call h.Create(rec, req)
// directly, and confirm rec.Code is http.StatusCreated and the decoded
// JSON body has the fields you sent.
func TestItemHandler_Create_ReturnsCreated(t *testing.T) {
	t.Skip("TODO: implement TestItemHandler_Create_ReturnsCreated")
}

// TestItemHandler_Create_ReturnsBadRequestOnInvalidQuantity should send
// an item with Quantity: 0 and confirm rec.Code is http.StatusBadRequest.
func TestItemHandler_Create_ReturnsBadRequestOnInvalidQuantity(t *testing.T) {
	t.Skip("TODO: implement TestItemHandler_Create_ReturnsBadRequestOnInvalidQuantity")
}
