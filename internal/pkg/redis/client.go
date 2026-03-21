// Package redis provides utilities for creating and configuring a Redis client.
// It loads connection settings from environment variables.
package redis

import (
	"crypto/tls"

	"github.com/redis/go-redis/v9"
)

// NewClient creates a new Redis client using configuration loaded from environment variables.
// The envPrefix parameter is used to namespace the environment variable names.
// It returns an error if the configuration cannot be loaded.
func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)

	if err != nil {
		return nil, err
	}

	opts := &redis.Options{
		Addr:     cfg.Address + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	if cfg.UseTLS {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return redis.NewClient(opts), nil
}
