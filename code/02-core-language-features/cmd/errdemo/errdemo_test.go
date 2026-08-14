package main

import (
	"strings"
	"testing"
)

// TestDivide is the "your first Go test" example from the lecture: a
// _test.go file, a TestX(t *testing.T) function, nothing else required.
func TestDivide(t *testing.T) {
	got, err := divide(10, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 5 {
		t.Errorf("divide(10, 2) = %v, want 5", got)
	}
}

// TestDivideByZero checks the failure path: divide should report the
// problem as an error, not panic, and the message should actually say
// something useful.
func TestDivideByZero(t *testing.T) {
	_, err := divide(10, 0)
	if err == nil {
		t.Fatal("expected an error dividing by zero, got nil")
	}
	if !strings.Contains(err.Error(), "zero") {
		t.Errorf("error message %q doesn't mention the actual problem", err.Error())
	}
}
