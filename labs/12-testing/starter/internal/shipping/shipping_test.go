package shipping

import "testing"

// Exercise 1: fill this in so it calls ExpressSurcharge(80), compares
// the result to 92, and reports a failure with t.Errorf if it doesn't
// match. Once it passes, temporarily change ExpressSurcharge in
// shipping.go to multiply by 1.10 instead of 1.15, run `go test`
// again, and read exactly what the failure output shows you (and
// doesn't show you). Put the 1.15 back afterwards.
func TestExpressSurcharge(t *testing.T) {
	t.Skip("TODO exercise 1: assert ExpressSurcharge(80) == 92")
}

// Exercise 2: CostPerParcel panics instead of returning an error when
// parcels is 0. Use defer/recover to write a test that fails if
// CostPerParcel(90, 0) *doesn't* panic. Once it passes, temporarily
// remove the panic from CostPerParcel in shipping.go and confirm this
// test correctly starts failing - that's how you know it was actually
// checking something.
func TestCostPerParcelZeroPanics(t *testing.T) {
	t.Skip("TODO exercise 2: recover from CostPerParcel(90, 0) and t.Error if recover() returns nil")
}

// Exercise 3: table-driven test. Build a slice of at least four cases
// (name, input weight, expected cost) covering ZoneRate's four
// branches, then loop over them with t.Run(tc.name, ...). Once it's
// passing, run:
//
//	go test -run TestZoneRate/<one case name, spaces as underscores>
//
// to confirm you can target a single case by name.
func TestZoneRate(t *testing.T) {
	t.Skip("TODO exercise 3: build a table of cases and t.Run each one")
}

// Exercise 5: benchmark ExpressSurcharge. Loop b.N times calling it,
// then run:
//
//	go test -bench=. -benchmem
//
// Note the ns/op and allocation counts. Try changing ExpressSurcharge's
// implementation slightly and see whether the numbers move.
func BenchmarkExpressSurcharge(b *testing.B) {
	b.Skip("TODO exercise 5: for i := 0; i < b.N; i++ { ExpressSurcharge(80) }")
}
