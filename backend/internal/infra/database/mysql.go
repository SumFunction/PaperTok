package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config holds the MySQL connection configuration.
type Config struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:         "localhost",
		Port:         3306,
		Username:     "root",
		Password:     "",
		Database:     "papertok",
		MaxOpenConns: 25,
		MaxIdleConns: 5,
		MaxLifetime:  5 * time.Minute,
	}
}

// MySQL implements the Connector interface for MySQL database.
type MySQL struct {
	db *sql.DB
	mu sync.RWMutex
}

// Ensure MySQL implements Connector interface.
var _ Connector = (*MySQL)(nil)

// New creates a new MySQL connection.
func New(cfg Config) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Verify connection is alive
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &MySQL{db: db}, nil
}

// DB returns the underlying database instance.
func (m *MySQL) DB() DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.db
}

// Close closes the database connection.
func (m *MySQL) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// Ping verifies the database connection is alive.
func (m *MySQL) Ping(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	return m.db.PingContext(ctx)
}

// BeginTx starts a new transaction.
func (m *MySQL) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.db.BeginTx(ctx, opts)
}

// WithTx executes a function within a transaction.
// The transaction is committed if the function returns nil,
// otherwise it is rolled back.
func (m *MySQL) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("exec error: %w, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
