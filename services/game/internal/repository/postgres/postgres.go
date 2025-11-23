package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Client wraps PostgreSQL client
type Client struct {
	db *sql.DB
}

// NewClient creates a new PostgreSQL client
func NewClient(connStr string, maxConns, maxIdle int) (*Client, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxIdle)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{db: db}, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.db.Close()
}

// GetDB returns the underlying database connection
func (c *Client) GetDB() *sql.DB {
	return c.db
}

