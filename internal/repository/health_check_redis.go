// Package repository provides data access layer implementations for the Bookmark Management API.
// Each repository manages a specific data source, such as a Redis instance.
package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// HealthCheckRedis defines the interface for checking Redis connectivity.
//
//go:generate mockery --name HealthCheckRedis --filename HealthCheckRedis.go
type HealthCheckRedis interface {
	// HealthCheck verifies that the Redis server is reachable.
	// It returns an error if the ping operation fails or times out.
	HealthCheck(ctx context.Context) error
}

// healthCheckRedis is the concrete implementation of HealthCheckRedis.
// It uses a Redis client to perform connectivity checks.
type healthCheckRedis struct {
	redisClient *redis.Client
}

// NewHealthCheckRedis creates a new HealthCheckRedis with the given Redis client.
func NewHealthCheckRedis(redisClient *redis.Client) HealthCheckRedis {
	return &healthCheckRedis{redisClient: redisClient}
}

// HealthCheck sends a PING command to the Redis server to verify connectivity.
// It applies a 2-second timeout to prevent indefinite blocking.
// It returns an error if the ping fails or the timeout is exceeded.
func (h *healthCheckRedis) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return h.redisClient.Ping(ctx).Err()
}
