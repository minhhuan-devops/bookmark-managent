package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/senn404/bookmark-managent/internal/api"
	"github.com/senn404/bookmark-managent/internal/config"
)

// TestPasswordEndpoint is an integration test that verifies the /gen-pass
// endpoint returns a successful response with the expected password length
// through the full API stack.
func TestPasswordEndpoint(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus int
		expectedLen    int
	}{
		{
			name: "success",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)
				return respRecorder
			},

			expectedStatus: http.StatusOK,
			expectedLen:    16, //Vì là json nên len dài hơn
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := api.New(&config.Config{}, nil)

			rec := tc.setupTestHTTP(app)

			var response struct {
				Password string `json:"password"`
			}
			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedLen, len(response.Password))
		})
	}
}
