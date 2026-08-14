// Command abstractfactory demonstrates Abstract Factory: producing a
// coordinated family of related objects (a connection + a query builder)
// without callers naming a concrete type anywhere.
package main

import "fmt"

type Connection interface {
	Connect() string
}

type QueryBuilder interface {
	Select(table string) string
}

// --- MySQL-flavored family ---

type mysqlConn struct{}

func (mysqlConn) Connect() string { return "connected via MySQL protocol" }

type mysqlQueryBuilder struct{}

func (mysqlQueryBuilder) Select(table string) string {
	return fmt.Sprintf("SELECT * FROM `%s` -- MySQL dialect", table)
}

// --- Postgres-flavored family ---

type pgConn struct{}

func (pgConn) Connect() string { return "connected via Postgres wire protocol" }

type pgQueryBuilder struct{}

func (pgQueryBuilder) Select(table string) string {
	return fmt.Sprintf(`SELECT * FROM "%s" -- Postgres dialect`, table)
}

// ConnectorFactory is the abstract factory itself: one interface that
// produces an entire matched family of products.
type ConnectorFactory interface {
	NewConnection() Connection
	NewQueryBuilder() QueryBuilder
}

type MySQLFactory struct{}

func (MySQLFactory) NewConnection() Connection     { return mysqlConn{} }
func (MySQLFactory) NewQueryBuilder() QueryBuilder { return mysqlQueryBuilder{} }

type PostgresFactory struct{}

func (PostgresFactory) NewConnection() Connection     { return pgConn{} }
func (PostgresFactory) NewQueryBuilder() QueryBuilder { return pgQueryBuilder{} }

// describeStack is written entirely against ConnectorFactory - it never
// mentions MySQLFactory or PostgresFactory by name.
func describeStack(f ConnectorFactory) {
	conn := f.NewConnection()
	qb := f.NewQueryBuilder()
	fmt.Println(conn.Connect())
	fmt.Println(qb.Select("guests"))
}

func main() {
	fmt.Println("-- MySQL stack --")
	describeStack(MySQLFactory{})

	fmt.Println("-- Postgres stack --")
	describeStack(PostgresFactory{})
}
