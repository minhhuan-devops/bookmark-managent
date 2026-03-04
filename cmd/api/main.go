package main

import (
	"github.com/joho/godotenv"
	"github.com/senn404/bookmark-managent/internal/api"
)

// @title Bookmark Management API
// @version 1.0
// @description This is a simple bookmark management API.
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg, err := api.NewConfig("")
	if err != nil {
		panic(err)
	}

	app := api.New(cfg)
	app.Start()
}
