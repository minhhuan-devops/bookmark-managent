// Package main is the entry point for the Bookmark Management API server.
// It initializes the application configuration and starts the HTTP server.
package main

import (
	"github.com/rs/zerolog/log"
	"github.com/senn404/bookmark-managent/internal/api"
	"github.com/senn404/bookmark-managent/internal/config"
	"github.com/senn404/bookmark-managent/internal/pkg/logger"
	"github.com/senn404/bookmark-managent/internal/pkg/redis"
)

// @title Bookmark Management API
// @version 1.0
// @description This is a simple bookmark management API.
// @host localhost:8080
// @BasePath /
// @schemes http
//
// main initializes the application configuration from environment variables
// and starts the Bookmark Management API server. It panics if the configuration
// cannot be loaded.
func main() {
	logger.SetLogLevel()

	cfg, err := config.NewConfig("")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	redisClient, err := redis.NewClient("")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to redis")
	}

	log.Info().Str("port", cfg.AppPort).Msg("Starting server")

	app := api.New(cfg, redisClient)
	if err := app.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
