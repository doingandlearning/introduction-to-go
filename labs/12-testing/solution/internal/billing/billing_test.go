package billing

import "testing"

// Exercise 1: a basic test. No assertion library - just an if
// statement and t.Errorf. (Break TenPercentOff to multiply by 0.8 and
// re-run `go test` to see the failure output for yourself; the
// message names the exact call and the got/want values, nothing more -
// it won't tell you *why* the numbers differ.)
func TestTenPercentOff(t *testing.T) {
	got := TenPercentOff(100)
	want := 90.0
	if got != want {
		t.Errorf("TenPercentOff(100) = %v, want %v", got, want)
	}
}

// A supporting test for Divide's ordinary (non-panicking) path.
func TestDivide(t *testing.T) {
	got := Divide(10, 2)
	want := 5.0
	if got != want {
		t.Errorf("Divide(10, 2) = %v, want %v", got, want)
	}
}

// Exercise 2: testing for a panic with defer/recover, wired up by hand
// - there's no assertPanics helper in the standard library. If you
// remove the panic from Divide, recover() returns nil and this test
// fails, which is how you confirm it was actually checking something.
func TestDivideByZeroPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected Divide(10, 0) to panic, but it didn't")
		}
	}()
	Divide(10, 0)
}

// Exercise 3: table-driven test with t.Run per case. This table
// deliberately does NOT cover TierDiscount's >= 100 branch (the 20%
// tier) - that gap is what exercise 4 finds with `go test -cover` and
// fills in below, in TestTierDiscountBigOrder.
//
// Run a single case by name with:
//
//	go test -run 'TestTierDiscount/boundary_at_fifty_counts_as_mid_tier'
func TestTierDiscount(t *testing.T) {
	cases := []struct {
		name  string
		price float64
		want  float64
	}{
		{"zero price stays zero", 0, 0},
		{"small price gets flat fiver floored at zero", 4, 0},
		{"mid price gets ten percent off", 60, 54},
		{"boundary at fifty counts as mid tier", 50, 45},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := TierDiscount(tc.price)
			if got != tc.want {
				t.Errorf("TierDiscount(%v) = %v, want %v", tc.price, got, tc.want)
			}
		})
	}
}

// Exercise 4: found via `go test -coverprofile=cover.out ./... && go
// tool cover -html=cover.out` - the >= 100 branch of TierDiscount (20%
// off for big orders) was never exercised by TestTierDiscount above.
// This test closes that gap.
func TestTierDiscountBigOrder(t *testing.T) {
	got := TierDiscount(150)
	want := 120.0
	if got != want {
		t.Errorf("TierDiscount(150) = %v, want %v", got, want)
	}
}

func TestApplyDiscount(t *testing.T) {
	got := ApplyDiscount(100, TenPercentOff)
	want := 90.0
	if got != want {
		t.Errorf("ApplyDiscount(100, TenPercentOff) = %v, want %v", got, want)
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
func BenchmarkTenPercentOff(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TenPercentOff(100)
	}
}
