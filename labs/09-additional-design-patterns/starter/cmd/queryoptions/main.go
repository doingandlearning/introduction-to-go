// Command queryoptions is Exercise 5: rewrite the ReportBuilder from
// Exercise 4 using functional options instead of chained methods, and
// try to replicate the "Columns and Table are required" validation
// using only functional options.
package main

import "fmt"

// TODO (Exercise 5): declare a ReportConfig struct with the same fields
// as ReportBuilder (columns, table, orderBy).

// TODO (Exercise 5): declare a ReportOption func(*ReportConfig) type.

// TODO (Exercise 5): implement WithColumns(cols string) ReportOption,
// WithTable(name string) ReportOption, and WithOrderBy(clause string)
// ReportOption, each returning a closure that sets the matching field.

// TODO (Exercise 5): implement NewReport(opts ...ReportOption)
// (string, error) that applies every option to a zero-value
// ReportConfig, then tries to validate and build the query string -
// same rules as Exercise 4's Build().

func main() {
	// TODO (Exercise 5a): call NewReport with WithColumns, WithTable, and
	// WithOrderBy, and print the result.

	// TODO (Exercise 5b): call NewReport with only WithColumns (no
	// WithTable) and confirm you still get an error back.

	// TODO (Exercise 5c): write a one-sentence comment here: did
	// replicating "Table is required" feel natural with functional
	// options, or did it feel like you were working against the pattern?

	fmt.Println("implement the TODOs above")
}
