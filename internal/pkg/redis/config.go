package redis

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// config holds the Redis connection settings.
// Fields are populated from environment variables with sensible defaults.
type config struct {
	// Address is the host and port of the Redis server (e.g., "localhost:6379").
	Address  string `default:"localhost:6379" envconfig:"REDIS_ADDRESS"`
	// Password is the authentication password for the Redis server.
	Password string `default:"" envconfig:"REDIS_PASSWORD"`
	// DB is the Redis database number to select after connecting.
	DB       int    `default:"0" envconfig:"REDIS_DB"`
}

// newConfig creates a new Redis config by loading a .env file (if present) and
// processing environment variables with the given prefix.
// It returns an error if the environment variables cannot be parsed.
func newConfig(envPrefix string) (*config, error) {
	_ = godotenv.Load()

	cfg := &config{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
