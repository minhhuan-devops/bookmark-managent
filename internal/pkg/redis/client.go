// Package redis provides utilities for creating and configuring a Redis client.
// It loads connection settings from environment variables.
package redis

import "github.com/redis/go-redis/v9"

// NewClient creates a new Redis client using configuration loaded from environment variables.
// The envPrefix parameter is used to namespace the environment variable names.
// It returns an error if the configuration cannot be loaded.
func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)

	if err != nil {
		return nil, err
	}

	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	}), nil
}
