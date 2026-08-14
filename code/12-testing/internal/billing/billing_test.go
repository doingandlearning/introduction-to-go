package billing

import "testing"

// A basic test: call the function, compare the result by hand. No
// assertion library, no framework magic - just an if statement and a
// call to t.Errorf when the values don't match.
func TestTenPercentOff(t *testing.T) {
	got := TenPercentOff(100)
	want := 90.0
	if got != want {
		t.Errorf("TenPercentOff(100) = %v, want %v", got, want)
	}
}

// A sanity-check test for the non-panicking branch of Divide.
func TestDivide(t *testing.T) {
	got := Divide(10, 2)
	want := 5.0
	if got != want {
		t.Errorf("Divide(10, 2) = %v, want %v", got, want)
	}
}

// Testing for a panic: there's no assertPanics helper in the standard
// library, so defer/recover is wired up by hand, inside the test
// itself. If recover() returns nil, the deferred function never saw a
// panic, so the test fails.
func TestDivideByZeroPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected Divide(10, 0) to panic, but it didn't")
		}
	}()
	Divide(10, 0)
}

// Table-driven test: ordinary Go data plus a loop, wrapped per case in
// t.Run so each case reports pass/fail independently and can be
// targeted by name with `go test -run`.
func TestTierDiscount(t *testing.T) {
	cases := []struct {
		name  string
		price float64
		want  float64
	}{
		{"zero price stays zero", 0, 0},
		{"small price gets flat fiver floored at zero", 4, 0},
		{"mid price gets ten percent off", 60, 54},
		{"large price gets twenty percent off", 150, 120},
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

// ApplyDiscount is tested with TenPercentOff plugged in as the
// strategy - proof that ApplyDiscount only cares about the function
// shape, not which specific discount it was handed.
func TestApplyDiscount(t *testing.T) {
	got := ApplyDiscount(100, TenPercentOff)
	want := 90.0
	if got != want {
		t.Errorf("ApplyDiscount(100, TenPercentOff) = %v, want %v", got, want)
	}
}

// BenchmarkTenPercentOff times TenPercentOff. go test decides how many
// times to run the loop body (b.N), adjusting it across runs until the
// timing is stable - that number is not something you choose.
func BenchmarkTenPercentOff(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TenPercentOff(100)
	}
}
