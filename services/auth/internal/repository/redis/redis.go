package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Connect creates a connection to Redis
func Connect(address, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	// Ping Redis to verify connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

