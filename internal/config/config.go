// Package config provides application configuration management for the
// Bookmark Management API. It loads configuration values from environment
// variables with sensible defaults using envconfig.
package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds the application configuration values. Fields are populated
// from environment variables, with fallback to the specified default values.
type Config struct {
	// AppPort is the port number the HTTP server listens on.
	AppPort string `default:"8080" envconfig:"APP_PORT"`
	// ServiceName is the name used to identify this service in health checks
	// and monitoring.
	ServiceName string `default:"bookmark_service" envconfig:"SERVICE_NAME"`
	// InstanceID is a unique identifier for this running instance. If left
	// empty, a UUID will be auto-generated at startup.
	InstanceID string `default:"" envconfig:"INSTANCE_ID"`
	// SwaggerHost is the host shown in Swagger UI
	// If empty, defaults to localhost:8080.
	SwaggerHost string `default:"" envconfig:"SWAGGER_HOST"`
}

// NewConfig creates a new Config by loading a .env file (if present) and
// processing environment variables with the given prefix. If envPrefix is
// empty, all matching environment variable names are used without a prefix.
// It returns an error if the environment variables cannot be parsed.
func NewConfig(envPrefix string) (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
