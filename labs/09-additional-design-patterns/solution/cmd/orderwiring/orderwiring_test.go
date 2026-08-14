// Exercise 7: prove TicketService behaves correctly when wired with the
// FakeInventory and FakeEmailer from Exercise 6 - no real warehouse or
// SMTP server involved.
package main

import "testing"

func TestBookWithFakes(t *testing.T) {
	fakeInv := &FakeInventory{}
	fakeMailer := &FakeEmailer{}
	svc := NewTicketService(fakeInv, fakeMailer)

	if err := svc.Book("Go Programming - Day 4", "delegate@example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fakeInv.reserved) != 1 || fakeInv.reserved[0] != "Go Programming - Day 4" {
		t.Errorf("fakeInv.reserved = %v, want [%q]", fakeInv.reserved, "Go Programming - Day 4")
	}

	if len(fakeMailer.sent) != 1 {
		t.Fatalf("fakeMailer.sent = %v, want exactly one message", fakeMailer.sent)
	}

	wantMessage := "delegate@example.com: your ticket for Go Programming - Day 4 is confirmed"
	if fakeMailer.sent[0] != wantMessage {
		t.Errorf("fakeMailer.sent[0] = %q, want %q", fakeMailer.sent[0], wantMessage)
	}
}
