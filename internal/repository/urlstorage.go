// Package repository provides data access layer implementations for the Bookmark Management API.
// Each repository manages a specific data source, such as a Redis instance.
package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// URLStorage defines the interface for storing and retrieving shortened URLs.
//
//go:generate mockery --name URLStorage --filename URLStorage.go
type URLStorage interface {
	// StoreURL stores the given URL under the specified code with a TTL of expTime.
	// It uses the NX (set-if-not-exists) mode to avoid overwriting existing codes.
	// It returns "OK" if the code was successfully stored, or an error on failure.
	StoreURL(ctx context.Context, code, url string, expTime time.Duration) (string, error)
}

// urlStorage is the concrete implementation of URLStorage.
// It persists URLs in a Redis instance using atomic set-if-not-exists operations.
type urlStorage struct {
	redisClient *redis.Client
}

// NewURLStorage creates a new URLStorage with the given Redis client.
func NewURLStorage(redisClent *redis.Client) URLStorage {
	return &urlStorage{
		redisClient: redisClent,
	}
}

// StoreURL stores the given url in Redis under code with the specified expiration.
// It uses the NX mode to ensure the key is only set if it does not already exist.
// It returns the Redis result string, or an error if the operation fails.
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, expTime time.Duration) (string, error) {
	return s.redisClient.SetArgs(ctx, code, url, redis.SetArgs{
		Mode: "NX",
		TTL:  expTime,
	}).Result()
}
