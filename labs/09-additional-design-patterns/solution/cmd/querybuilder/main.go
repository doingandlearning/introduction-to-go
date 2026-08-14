// Command querybuilder is Exercise 4: a chained Builder with validation
// that only succeeds once every required step has been called.
package main

import "fmt"

type ReportBuilder struct {
	columns string
	table   string
	orderBy string
}

func NewReportBuilder() *ReportBuilder {
	return &ReportBuilder{}
}

func (b *ReportBuilder) Columns(cols string) *ReportBuilder {
	b.columns = cols
	return b
}

func (b *ReportBuilder) Table(name string) *ReportBuilder {
	b.table = name
	return b
}

func (b *ReportBuilder) OrderBy(clause string) *ReportBuilder {
	b.orderBy = clause
	return b
}

func (b *ReportBuilder) Build() (string, error) {
	if b.columns == "" || b.table == "" {
		return "", fmt.Errorf("columns and table are required")
	}
	query := fmt.Sprintf("SELECT %s FROM %s", b.columns, b.table)
	if b.orderBy != "" {
		query += " ORDER BY " + b.orderBy
	}
	return query, nil
}

func main() {
	q, err := NewReportBuilder().
		Columns("event, tickets_sold").
		Table("sales_report").
		OrderBy("tickets_sold DESC").
		Build()
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(q)
	}

	_, err = NewReportBuilder().Columns("event").Build()
	fmt.Println("expected error:", err)
}
