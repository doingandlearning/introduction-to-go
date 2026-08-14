// Exercise 7: prove the Part A Strategy functions (and swapping between
// them) behave the way the lecture claimed - trivially, with no mocking.
package main

import "testing"

func TestDiscountFunctions(t *testing.T) {
	price := 200.0

	cases := []struct {
		name string
		fn   TicketDiscount
		want float64
	}{
		{"no discount", NoDiscount, 200.0},
		{"member discount", MemberDiscount, 170.0},
		{"early bird discount", EarlyBirdDiscount, 140.0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.fn(price)
			if got != c.want {
				t.Errorf("%s(%.2f) = %.2f, want %.2f", c.name, price, got, c.want)
			}
		})
	}
}

func TestApplyTicketDiscount(t *testing.T) {
	price := 200.0

	got := ApplyTicketDiscount(price, NoDiscount)
	if got != 200.0 {
		t.Errorf("ApplyTicketDiscount(price, NoDiscount) = %.2f, want 200.00", got)
	}

	got = ApplyTicketDiscount(price, MemberDiscount)
	if got != 170.0 {
		t.Errorf("ApplyTicketDiscount(price, MemberDiscount) = %.2f, want 170.00", got)
	}

	// Swapping the strategy argument must change the result - that's
	// the entire point of depending on a function type.
	if ApplyTicketDiscount(price, NoDiscount) == ApplyTicketDiscount(price, EarlyBirdDiscount) {
		t.Errorf("expected swapping the strategy from NoDiscount to EarlyBirdDiscount to change the result")
	}
}
