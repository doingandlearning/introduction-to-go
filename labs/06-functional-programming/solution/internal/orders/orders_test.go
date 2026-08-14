package orders

import "testing"

func TestNewCoffeeOrder(t *testing.T) {
	// Zero options: defaults must survive untouched.
	plain := NewCoffeeOrder()
	if plain.Size() != "medium" {
		t.Errorf("zero options: Size() = %q, want %q", plain.Size(), "medium")
	}
	if plain.ExtraShot() {
		t.Errorf("zero options: ExtraShot() = true, want false")
	}
	if plain.OatMilk() {
		t.Errorf("zero options: OatMilk() = true, want false")
	}

	// One option: only the overridden field should change.
	large := NewCoffeeOrder(WithSize("large"))
	if large.Size() != "large" {
		t.Errorf("one option: Size() = %q, want %q", large.Size(), "large")
	}
	if large.ExtraShot() {
		t.Errorf("one option: ExtraShot() = true, want false")
	}
	if large.OatMilk() {
		t.Errorf("one option: OatMilk() = true, want false")
	}

	// Multiple options stacked: every override should apply together.
	loaded := NewCoffeeOrder(WithSize("small"), WithExtraShot(), WithOatMilk())
	if loaded.Size() != "small" {
		t.Errorf("stacked options: Size() = %q, want %q", loaded.Size(), "small")
	}
	if !loaded.ExtraShot() {
		t.Errorf("stacked options: ExtraShot() = false, want true")
	}
	if !loaded.OatMilk() {
		t.Errorf("stacked options: OatMilk() = false, want true")
	}
}
