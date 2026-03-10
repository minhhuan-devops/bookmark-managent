package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/redis/go-redis/v9"
	"github.com/senn404/bookmark-managent/internal/api"
	"github.com/senn404/bookmark-managent/internal/config"
	pkgRedis "github.com/senn404/bookmark-managent/internal/pkg/redis"
)

// TestPasswordEndpoint is an integration test that verifies the /gen-pass
// endpoint returns a successful response with the expected password length
// through the full API stack.

type errorResponse struct {
	Error string `json:"error"`
}

type shortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TestShortenURLEndpoint(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		setupTestHTTP  func(api api.Engine) *httptest.ResponseRecorder
		setupMockRedis func() *redis.Client

		expectedStatus int
		expectedRespon shortenURLResponse
	}{
		{
			name: "success",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {

				body, _ := json.Marshal(map[string]any{
					"url":      "huanops.com",
					"exp_time": 60,
				})

				req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			setupMockRedis: func() *redis.Client {
				return pkgRedis.InitMockRedis(t)
			},

			expectedStatus: http.StatusOK,
			expectedRespon: shortenURLResponse{
				Code:    "",
				Message: "Shorten URL generated successfully!",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redisClient := tc.setupMockRedis()

			app := api.New(&config.Config{}, redisClient)

			rec := tc.setupTestHTTP(app)

			respon := shortenURLResponse{}
			json.Unmarshal(rec.Body.Bytes(), &respon)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedRespon.Message, respon.Message)
		})
	}
}
