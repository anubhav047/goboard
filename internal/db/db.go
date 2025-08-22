package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates a new database connection pool.
func Connect(connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping the database to verify the connection is live.
	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close() // Close the pool if ping fails.
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
