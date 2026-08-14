// Command orderwiring is Exercise 6: build a small Service that depends
// on two interfaces, wire it by hand with real implementations, then
// wire it again with fakes, and confirm the Service code never changes.
package main

import "fmt"

type Inventory interface {
	Reserve(item string) error
}

type Emailer interface {
	Send(to, message string)
}

type TicketService struct {
	inv    Inventory
	mailer Emailer
}

func NewTicketService(inv Inventory, mailer Emailer) *TicketService {
	return &TicketService{inv: inv, mailer: mailer}
}

func (s *TicketService) Book(item, email string) error {
	if err := s.inv.Reserve(item); err != nil {
		return err
	}
	s.mailer.Send(email, "your ticket for "+item+" is confirmed")
	return nil
}

// --- real implementations ---

type WarehouseInventory struct{}

func (WarehouseInventory) Reserve(item string) error {
	fmt.Println("[warehouse] reserved:", item)
	return nil
}

type SMTPEmailer struct{}

func (SMTPEmailer) Send(to, message string) {
	fmt.Printf("[smtp] to=%s message=%q\n", to, message)
}

// --- fake implementations ---

type FakeInventory struct {
	reserved []string
}

func (i *FakeInventory) Reserve(item string) error {
	i.reserved = append(i.reserved, item)
	return nil
}

type FakeEmailer struct {
	sent []string
}

func (e *FakeEmailer) Send(to, message string) {
	e.sent = append(e.sent, to+": "+message)
}

func main() {
	fmt.Println("-- real wiring --")
	realService := NewTicketService(WarehouseInventory{}, SMTPEmailer{})
	if err := realService.Book("Go Programming - Day 4", "delegate@example.com"); err != nil {
		fmt.Println("booking failed:", err)
	}

	fmt.Println("-- fake wiring (what a test would do) --")
	fakeInv := &FakeInventory{}
	fakeMailer := &FakeEmailer{}
	testService := NewTicketService(fakeInv, fakeMailer)
	if err := testService.Book("Go Programming - Day 4", "delegate@example.com"); err != nil {
		fmt.Println("booking failed:", err)
	}

	fmt.Println("fake inventory reserved:", fakeInv.reserved)
	fmt.Println("fake mailer sent:", fakeMailer.sent)

	// TicketService, NewTicketService, and Book are byte-for-byte
	// identical in both wirings above - only the two arguments passed
	// into NewTicketService changed. That's the entire value of
	// depending on interfaces instead of concrete types.
}
