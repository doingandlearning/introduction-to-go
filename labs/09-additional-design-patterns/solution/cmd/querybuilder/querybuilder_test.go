// Exercise 7: prove ReportBuilder's Build() both succeeds on a fully
// chained call and returns an error when a required step is missing.
package main

import "testing"

func TestBuildSuccess(t *testing.T) {
	got, err := NewReportBuilder().
		Columns("event, tickets_sold").
		Table("sales_report").
		OrderBy("tickets_sold DESC").
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "SELECT event, tickets_sold FROM sales_report ORDER BY tickets_sold DESC"
	if got != want {
		t.Errorf("Build() = %q, want %q", got, want)
	}
}

func TestBuildMissingTable(t *testing.T) {
	got, err := NewReportBuilder().Columns("event").Build()
	if err == nil {
		t.Fatalf("expected an error when Table() was never called, got nil result %q", got)
	}
}
