package mathutils

import "testing"

func TestSafeDivide(t *testing.T) {
	got, err := SafeDivide(10, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 5 {
		t.Errorf("SafeDivide(10, 2) = %v, want 5", got)
	}

	_, err = SafeDivide(10, 0)
	if err == nil {
		t.Fatal("expected an error dividing by zero, got nil")
	}
}

func TestAdd(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{"a few positives", []int{1, 2, 3}, 6},
		{"single value", []int{7}, 7},
		{"no arguments at all", []int{}, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Add(tc.nums...)
			if got != tc.want {
				t.Errorf("Add(%v) = %v, want %v", tc.nums, got, tc.want)
			}
		})
	}
}
