// This test file is already complete — you're not writing it. It's the
// specification for Exercise 2: run `go test ./...` now, before touching
// decisions.go, and both tests fail. Implement LateFeeTier and
// DeskSchedule until they pass. Writing a test like this yourself is
// Topic 12's job, not this one's.
package catalog

import "testing"

func TestLateFeeTier(t *testing.T) {
	tests := []struct {
		name     string
		daysLate int
		want     string
	}{
		{"well within grace period", 0, "none"},
		{"just under warning threshold", 6, "none"},
		{"start of warning window", 7, "warning"},
		{"end of warning window", 29, "warning"},
		{"start of suspension", 30, "suspended"},
		{"far overdue", 90, "suspended"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LateFeeTier(tt.daysLate)
			if got != tt.want {
				t.Errorf("LateFeeTier(%d) = %q, want %q", tt.daysLate, got, tt.want)
			}
		})
	}
}

func TestDeskSchedule(t *testing.T) {
	tests := []struct {
		name string
		day  int
		want string
	}{
		{"weekday", 3, "open 9am-6pm"},
		{"Saturday falls through", 6, "open 10am-2pm only"},
		{"Sunday", 7, "open 10am-2pm only"},
		{"unknown day", 0, "unknown day"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeskSchedule(tt.day)
			if got != tt.want {
				t.Errorf("DeskSchedule(%d) = %q, want %q", tt.day, got, tt.want)
			}
		})
	}

	// The fallthrough from Saturday's case into Sunday's is the whole
	// point of Exercise 2 — assert it directly, not just via the table
	// above, so breaking it produces an obvious, named failure.
	sat, sun := DeskSchedule(6), DeskSchedule(7)
	if sat != sun {
		t.Errorf("expected Saturday and Sunday to match via fallthrough, got %q and %q", sat, sun)
	}
}
