// Command builder demonstrates the Builder pattern: step-by-step
// assembly with chained calls, plus validation that only makes sense
// once every step is complete.
package main

import "fmt"

type QueryBuilder struct {
	cols  string
	table string
	where string
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (q *QueryBuilder) Select(cols string) *QueryBuilder {
	q.cols = cols
	return q
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
	q.table = table
	return q
}

func (q *QueryBuilder) Where(cond string) *QueryBuilder {
	q.where = cond
	return q
}

// Build validates the accumulated state, something a plain struct
// literal has no natural hook for.
func (q *QueryBuilder) Build() (string, error) {
	if q.cols == "" || q.table == "" {
		return "", fmt.Errorf("select and from are required")
	}
	if q.where == "" {
		return fmt.Sprintf("SELECT %s FROM %s", q.cols, q.table), nil
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", q.cols, q.table, q.where), nil
}

func main() {
	sql, err := NewQueryBuilder().
		Select("id, name").
		From("guests").
		Where("active = true").
		Build()
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(sql)
	}

	// Missing From() - Build() refuses before a malformed query ever runs.
	_, err = NewQueryBuilder().Select("id").Build()
	fmt.Println("expected error:", err)
}
