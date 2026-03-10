// Package api provides the HTTP server engine for the Bookmark Management API.
// It sets up the Gin router, registers all endpoint handlers, and manages
// the application lifecycle including Swagger documentation.
package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	_ "github.com/senn404/bookmark-managent/docs"
	"github.com/senn404/bookmark-managent/internal/config"
	"github.com/senn404/bookmark-managent/internal/handler"
	"github.com/senn404/bookmark-managent/internal/repository"
	"github.com/senn404/bookmark-managent/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine defines the interface for the API server. It provides methods
// to start the server and to serve HTTP requests directly, which is
// useful for testing.
type Engine interface {
	// Start starts the HTTP server and begins listening for incoming requests.
	Start() error
	// ServeHTTP implements the http.Handler interface, allowing the engine
	// to be used directly in tests with httptest.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// api is the concrete implementation of the Engine interface. It wraps
// a Gin engine and the application configuration.
type api struct {
	app         *gin.Engine
	cfg         *config.Config
	redisClient *redis.Client
}

// Start starts the HTTP server on the port specified in the application configuration.
// It returns an error if the server fails to start.
func (a *api) Start() error {
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

// ServeHTTP delegates HTTP request handling to the underlying Gin engine.
// This method satisfies the http.Handler interface and is primarily used
// for integration testing with httptest.
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

// registerEP registers all API endpoint handlers and Swagger documentation route.
// It initializes the service layer and handler layer using dependency injection,
// and maps each handler to its corresponding HTTP route.
func (a *api) registerEP() {
	// Create Handler
	//passwork handler
	passSvc := service.NewPassword()
	passHandler := handler.NewPasswordHandler(passSvc)
	//healthcheck handler
	healthCheckRedis := repository.NewHealthCheckRedis(a.redisClient)
	healthCheck := service.NewHealthCheck(a.cfg, healthCheckRedis)
	healthCheckHandler := handler.NewHealthCheckHandler(healthCheck)
	//url shorten handler
	urlStorage := repository.NewURLStorage(a.redisClient)
	urlService := service.NewShortenURLService(urlStorage)
	urlHandler := handler.NewShortenURLHandler(urlService)

	a.app.GET("/gen-pass", passHandler.GenPass)
	a.app.GET("/health-check", healthCheckHandler.HealthCheck)
	a.app.POST("/shorten", urlHandler.ShortenURL)

	//URL Storage

	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// New creates a new Engine instance with the given configuration.
// It initializes the Gin router with default middleware (logger and recovery),
// registers all endpoint handlers, and returns the Engine ready to start.
func New(cfg *config.Config, redisClient *redis.Client) Engine {
	a := &api{
		app:         gin.Default(),
		cfg:         cfg,
		redisClient: redisClient,
	}
	a.registerEP()
	return a
}
