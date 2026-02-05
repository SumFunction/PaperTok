package database

import (
	"context"
	"database/sql"
)

// DB defines the interface for database operations.
// This abstraction allows for different database implementations
// and facilitates testing with mock implementations.
type DB interface {
	// QueryContext executes a query that returns rows.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// ExecContext executes a query without returning any rows.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// BeginTx starts a transaction with the given options.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// Close closes the database connection.
	Close() error

	// PingContext verifies a connection to the database is still alive.
	PingContext(ctx context.Context) error
}

// Connector defines the interface for database connection management.
type Connector interface {
	// DB returns the database instance.
	DB() DB

	// Close closes the database connection.
	Close() error

	// Ping verifies the database connection is alive.
	Ping(ctx context.Context) error
}

// Executor defines a common interface for both DB and Tx.
// This allows functions to work with either a database connection or a transaction.
type Executor interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
