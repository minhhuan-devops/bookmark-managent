// Package handler provides HTTP request handlers for the Bookmark Management
// API. Each handler is responsible for parsing incoming requests, delegating
// business logic to the service layer, and returning appropriate HTTP responses.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/senn404/bookmark-managent/internal/service"
)

// healthCheckHandler is the concrete implementation of HealthCheckHandler.
// It delegates health status retrieval to the HealthCheck service.
type healthCheckHandler struct {
	svc service.HealthCheck
}

// HealthCheckHandler defines the interface for handling health check HTTP requests.
type HealthCheckHandler interface {
	// HealthCheck handles GET /health-check requests and returns the current
	// service health status.
	HealthCheck(c *gin.Context)
}

// NewHealthCheckHandler creates a new HealthCheckHandler with the given
// HealthCheck service. It uses dependency injection to allow easy testing
// with mock services.
func NewHealthCheckHandler(svc service.HealthCheck) HealthCheckHandler {
	return &healthCheckHandler{svc: svc}
}

// HealthCheck handles the GET /health-check endpoint.
// It retrieves the current service health status and returns it as a JSON response.
//
// @Summary Health Check
// @Tags health-check
// @Produce json
// @Success 200 {object} string
// @Failure 500 {object} errorResponse
// @Router /health-check [get]
func (h *healthCheckHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, h.svc.GetStatus())
}
