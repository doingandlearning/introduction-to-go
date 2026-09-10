package shipping

import "testing"

// Exercise 1: a basic test. No assertion library - just an if
// statement and t.Errorf.
func TestExpressSurcharge(t *testing.T) {
	got := ExpressSurcharge(80)
	want := 92.0
	if got != want {
		t.Errorf("ExpressSurcharge(80) = %v, want %v", got, want)
	}
}

// A supporting test for CostPerParcel's ordinary (non-panicking) path.
func TestCostPerParcel(t *testing.T) {
	got := CostPerParcel(90, 3)
	want := 30.0
	if got != want {
		t.Errorf("CostPerParcel(90, 3) = %v, want %v", got, want)
	}
}

// Exercise 2: testing for a panic with defer/recover, wired up by hand
// - there's no assertPanics helper in the standard library. If you
// remove the panic from CostPerParcel, recover() returns nil and this
// test fails, which is how you confirm it was actually checking
// something.
func TestCostPerParcelZeroPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected CostPerParcel(90, 0) to panic, but it didn't")
		}
	}()
	CostPerParcel(90, 0)
}

// Exercise 3: table-driven test with t.Run per case. This table
// deliberately does NOT cover ZoneRate's >= 20 branch (the heavy
// parcel rate) - that gap is what exercise 4 finds with `go test
// -cover` and fills in below, in TestZoneRateHeavyParcel.
//
// Run a single case by name with:
//
//	go test -run 'TestZoneRate/boundary_at_five_counts_as_standard_tier'
func TestZoneRate(t *testing.T) {
	cases := []struct {
		name   string
		weight float64
		want   float64
	}{
		{"zero weight stays zero", 0, 0},
		{"tiny parcel floored at zero", 0.5, 0},
		{"standard parcel rate applies", 8, 40},
		{"boundary at five counts as standard tier", 5, 25},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ZoneRate(tc.weight)
			if got != tc.want {
				t.Errorf("ZoneRate(%v) = %v, want %v", tc.weight, got, tc.want)
			}
		})
	}
}

// Exercise 4: found via `go test -coverprofile=cover.out ./... && go
// tool cover -html=cover.out` - the >= 20 branch of ZoneRate (the
// heavy parcel rate) was never exercised by TestZoneRate above. This
// test closes that gap.
func TestZoneRateHeavyParcel(t *testing.T) {
	got := ZoneRate(25)
	want := 100.0
	if got != want {
		t.Errorf("ZoneRate(25) = %v, want %v", got, want)
	}
}

func TestApplyRate(t *testing.T) {
	got := ApplyRate(80, ExpressSurcharge)
	want := 92.0
	if got != want {
		t.Errorf("ApplyRate(80, ExpressSurcharge) = %v, want %v", got, want)
	}
}

// Exercise 5: benchmark. go test decides how many times to run the
// loop body (b.N), adjusting it across runs until the timing settles -
// that number is not something you choose, the same way a stopwatch
// operator decides to "time a line for a second's worth of widgets,"
// however many that turns out to be.
//
// Run with:
//
//	go test -bench=. -benchmem
func BenchmarkExpressSurcharge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExpressSurcharge(80)
	}
}
