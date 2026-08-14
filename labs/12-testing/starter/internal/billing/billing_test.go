package billing

import "testing"

// Exercise 1: fill this in so it calls TenPercentOff(100), compares the
// result to 90, and reports a failure with t.Errorf if it doesn't
// match. Once it passes, temporarily change TenPercentOff in billing.go
// to multiply by 0.8 instead of 0.9, run `go test` again, and read
// exactly what the failure output shows you (and doesn't show you).
// Put the 0.9 back afterwards.
func TestTenPercentOff(t *testing.T) {
	t.Skip("TODO exercise 1: assert TenPercentOff(100) == 90")
}

// Exercise 2: Divide panics instead of returning an error when b is 0.
// Use defer/recover to write a test that fails if Divide(10, 0)
// *doesn't* panic. Once it passes, temporarily remove the panic from
// Divide in billing.go and confirm this test correctly starts failing
// - that's how you know it was actually checking something.
func TestDivideByZeroPanics(t *testing.T) {
	t.Skip("TODO exercise 2: recover from Divide(10, 0) and t.Error if recover() returns nil")
}

// Exercise 3: table-driven test. Build a slice of at least four cases
// (name, input price, expected result) covering TierDiscount's four
// branches, then loop over them with t.Run(tc.name, ...). Once it's
// passing, run:
//
//	go test -run TestTierDiscount/<one case name, spaces as underscores>
//
// to confirm you can target a single case by name.
func TestTierDiscount(t *testing.T) {
	t.Skip("TODO exercise 3: build a table of cases and t.Run each one")
}

// Exercise 5: benchmark TenPercentOff. Loop b.N times calling it, then
// run:
//
//	go test -bench=. -benchmem
//
// Note the ns/op and allocation counts. Try changing TenPercentOff's
// implementation slightly and see whether the numbers move.
func BenchmarkTenPercentOff(b *testing.B) {
	b.Skip("TODO exercise 5: for i := 0; i < b.N; i++ { TenPercentOff(100) }")
}
