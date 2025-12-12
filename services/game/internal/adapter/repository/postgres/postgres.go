package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Client struct {
	db *sql.DB
}

func NewClient(connStr string, maxConns, maxIdle int) (*Client, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxIdle)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) GetDB() *sql.DB {
	return c.db
}

func (c *Client) Ping() error {
	return c.db.Ping()
}

