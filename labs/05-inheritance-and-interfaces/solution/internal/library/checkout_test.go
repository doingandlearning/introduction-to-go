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
