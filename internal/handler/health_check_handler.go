// Package handler provides HTTP request handlers for the Bookmark Management API.
// Each handler is responsible for parsing incoming requests.
// Business logic is delegated to the service layer.
// Each handler returns an appropriate HTTP response to the client.
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
	// HealthCheck handles GET /health-check requests.
	// It returns the current service health status as a JSON response.
	HealthCheck(c *gin.Context)
}

// NewHealthCheckHandler creates a new HealthCheckHandler with the given HealthCheck service.
// It uses dependency injection to allow easy testing with mock services.
func NewHealthCheckHandler(svc service.HealthCheck) HealthCheckHandler {
	return &healthCheckHandler{svc: svc}
}

// HealthCheck handles the GET /health-check endpoint.
// It retrieves the current service health status.
// It returns the status as a JSON response.
//
// @Summary Health Check
// @Tags health-check
// @Produce json
// @Success 200 {object} service.HealthStatus
// @Failure 500 {object} errorResponse
// @Router /health-check [get]
func (h *healthCheckHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, h.svc.GetStatus(c))
}
