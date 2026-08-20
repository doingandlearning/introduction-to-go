// This test file is already complete — you're not writing it. It's the
// specification for Exercise 5: run `go test ./...` now, before
// touching checkout.go, and it fails on its very first assertion.
// Fix CheckOut until it passes. Writing a test like this yourself is
// Topic 12's job, not this one's.
package library

import "testing"

// TestCheckOutNilInterfaceGotcha confirms CheckOut returns a genuinely
// nil error on the success path, and a usable error once the patron is
// over their book limit (Exercise 5).
func TestCheckOutNilInterfaceGotcha(t *testing.T) {
	err := CheckOut("The Go Programming Language", 2, 5)
	if err != nil {
		t.Fatalf("CheckOut under the limit: err = %v, want nil", err)
	}

	err = CheckOut("The Go Programming Language", 5, 5)
	if err == nil {
		t.Fatal("CheckOut at the limit: want a non-nil error, got nil")
	}

	want := "could not check out: The Go Programming Language"
	if err.Error() != want {
		t.Errorf("err.Error() = %q, want %q", err.Error(), want)
	}
}
