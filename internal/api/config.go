package api

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort string `default:"8080" envconfig:"APP_PORT"`
}

func NewConfig(envPrefix string) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
