// Command queryoptions is Exercise 5: rewrite the ReportBuilder from
// Exercise 4 using functional options instead of chained methods, and
// try to replicate the "Columns and Table are required" validation
// using only functional options.
package main

import "fmt"

type ReportConfig struct {
	columns string
	table   string
	orderBy string
}

type ReportOption func(*ReportConfig)

func WithColumns(cols string) ReportOption {
	return func(c *ReportConfig) { c.columns = cols }
}

func WithTable(name string) ReportOption {
	return func(c *ReportConfig) { c.table = name }
}

func WithOrderBy(clause string) ReportOption {
	return func(c *ReportConfig) { c.orderBy = clause }
}

func NewReport(opts ...ReportOption) (string, error) {
	cfg := &ReportConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// The validation itself still works fine - the awkward part is what
	// comes before it: functional options give you no way to require
	// that WithColumns ran before WithTable, or that either ran at all,
	// until you're all the way down here checking the fields by hand.
	if cfg.columns == "" || cfg.table == "" {
		return "", fmt.Errorf("columns and table are required")
	}
	query := fmt.Sprintf("SELECT %s FROM %s", cfg.columns, cfg.table)
	if cfg.orderBy != "" {
		query += " ORDER BY " + cfg.orderBy
	}
	return query, nil
}

func main() {
	q, err := NewReport(
		WithColumns("event, tickets_sold"),
		WithTable("sales_report"),
		WithOrderBy("tickets_sold DESC"),
	)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(q)
	}

	_, err = NewReport(WithColumns("event"))
	fmt.Println("expected error:", err)

	// Exercise 5c: the validation works, but it's a post-hoc check
	// against a config struct, not a property of the type at each step.
	// Builder made "Table is required before Build succeeds" a
	// structural fact you couldn't avoid noticing while chaining -
	// functional options make it just another if-check at the bottom.
	// This felt forced: functional options are the right tool when
	// options are genuinely independent and optional, not when one is
	// mandatory or must precede another.
}
