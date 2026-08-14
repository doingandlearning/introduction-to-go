// Command querybuilder is Exercise 4: a chained Builder with validation
// that only succeeds once every required step has been called.
package main

import "fmt"

// TODO (Exercise 4): declare a ReportBuilder struct with fields for
// columns, table, and an optional orderBy clause.

// TODO (Exercise 4): implement NewReportBuilder() *ReportBuilder.

// TODO (Exercise 4): implement chained methods Columns(cols string),
// Table(name string), and OrderBy(clause string), each returning
// *ReportBuilder so calls can be chained.

// TODO (Exercise 4): implement Build() (string, error). It must return
// an error if Columns or Table was never called. If OrderBy was never
// called, omit the ORDER BY clause from the output.

func main() {
	// TODO (Exercise 4a): build a valid query - Columns, Table, and
	// OrderBy all set - and print the result.

	// TODO (Exercise 4b): build a query missing Table() and confirm
	// Build() returns a non-nil error. Print it.

	fmt.Println("implement the TODOs above")
}
