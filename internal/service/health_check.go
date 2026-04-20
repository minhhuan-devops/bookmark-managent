// Package service provides the business logic layer for the Bookmark
// Management API. It defines service interfaces and their implementations,
// separating business logic from HTTP handling concerns.
package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/senn404/bookmark-managent/internal/config"
	"github.com/senn404/bookmark-managent/internal/repository"
)

// HealthStatus represents the health status of the service, including
// identifying information about the running instance.
type HealthStatus struct {
	// Message indicates the overall health state (e.g., "OK").
	Message string `json:"message"`
	// ServiceName is the name of the service as configured.
	ServiceName string `json:"service_name"`
	// InstanceId is the unique identifier of this running instance.
	InstanceId string `json:"instance_id"`
}

// healthCheck is the concrete implementation of the HealthCheck interface.
// It holds a pre-computed HealthStatus that is returned on every status check.
type healthCheck struct {
	HealthCheckRedis repository.HealthCheckRedis
	status           HealthStatus
}

// HealthCheck defines the interface for retrieving the service health status.
//
//go:generate mockery --name HealthCheck --filename health_check.go
type HealthCheck interface {
	// GetStatus returns the current health status of the service.
	GetStatus(ctx context.Context) (status HealthStatus)
}

// NewHealthCheck creates a new HealthCheck service with the given configuration.
// If the InstanceID in the configuration is empty, a new UUID is automatically
// generated to uniquely identify this instance.
func NewHealthCheck(cfg *config.Config, htc repository.HealthCheckRedis) HealthCheck {
	instancID := cfg.InstanceID
	if instancID == "" {
		instancID = uuid.New().String()
	}

	return &healthCheck{
		htc,

		HealthStatus{
			Message:     "OK",
			ServiceName: cfg.ServiceName,
			InstanceId:  instancID,
		},
	}
}

// GetStatus returns the pre-computed health status of the service,
// including the service name, instance ID, and health message.
func (h *healthCheck) GetStatus(ctx context.Context) (status HealthStatus) {
	status = h.status
	if err := h.HealthCheckRedis.HealthCheck(ctx); err != nil {
		log.Error().Err(err).Msg("Redis health check failed")
		status.Message = "Internal Server Error"
	}
	return status
}
