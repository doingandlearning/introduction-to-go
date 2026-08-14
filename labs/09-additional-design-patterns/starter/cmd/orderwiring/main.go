// Command orderwiring is Exercise 6: build a small Service that depends
// on two interfaces, wire it by hand with real implementations, then
// wire it again with fakes, and confirm the Service code never changes.
package main

import "fmt"

// TODO (Exercise 6a): declare an Inventory interface with one method,
// Reserve(item string) error.

// TODO (Exercise 6a): declare an Emailer interface with one method,
// Send(to, message string).

// TODO (Exercise 6a): declare a TicketService struct holding an
// Inventory and an Emailer, with a constructor NewTicketService(inv
// Inventory, mailer Emailer) *TicketService.

// TODO (Exercise 6a): implement (*TicketService).Book(item, email
// string) error - call Reserve, and if it succeeds, call Send with a
// confirmation message. Return any error from Reserve.

// TODO (Exercise 6b): implement real types WarehouseInventory
// (Reserve always succeeds, logs what it reserved) and SMTPEmailer
// (Send just prints, pretending to email).

// TODO (Exercise 6c): implement fake types FakeInventory (records what
// was reserved in a slice, Reserve always succeeds) and FakeEmailer
// (records sent messages in a slice).

func main() {
	// TODO (Exercise 6d): wire a TicketService with the real
	// implementations in main, call Book once, and print the outcome.

	// TODO (Exercise 6e): wire a second TicketService with the fake
	// implementations, call Book once, and print what the fakes
	// recorded - confirming TicketService's code is identical in both
	// wirings.

	fmt.Println("implement the TODOs above")
}
