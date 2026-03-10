// Package endpoint provides integration tests for the Bookmark Management API.
// These tests verify end-to-end behavior by spinning up the full API stack
// (router, handlers, services) and making HTTP requests against it.
package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/senn404/bookmark-managent/internal/api"
	"github.com/senn404/bookmark-managent/internal/config"
	pkgRedis "github.com/senn404/bookmark-managent/internal/pkg/redis"
	"github.com/senn404/bookmark-managent/internal/service"
)

// TestHealthCheckEndpoint is an integration test that verifies the /health-check
// endpoint returns the correct health status, service name, and a valid UUID
// instance ID through the full API stack.
func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		setupTestHTTP  func(api api.Engine) *httptest.ResponseRecorder
		setupMockRedis func() *redis.Client

		expectedStatus      int
		expectedMessage     string
		expectedServiceName string
	}{
		{
			name: "success",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},
			setupMockRedis: func() *redis.Client {
				return pkgRedis.InitMockRedis(t)
			},

			expectedStatus:      http.StatusOK,
			expectedMessage:     "OK",
			expectedServiceName: "bookmark_service",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := config.Config{
				AppPort:     "8080",
				ServiceName: "bookmark_service",
				InstanceID:  "7dee4e26-66fd-44c1-a135-fe4ce9e4b8aa",
			}
			redisClient := tc.setupMockRedis()

			app := api.New(&cfg, redisClient)

			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var body service.HealthStatus
			json.Unmarshal(rec.Body.Bytes(), &body)
			assert.Equal(t, tc.expectedMessage, body.Message)
			assert.Equal(t, tc.expectedServiceName, body.ServiceName)

			// Sử dụng parse để test instance id có trả về đúng format và không rỗng khong
			_, err := uuid.Parse(body.InstanceId)
			assert.Equal(t, nil, err)
			assert.NotEqual(t, "", body.InstanceId)
		})
	}
}
