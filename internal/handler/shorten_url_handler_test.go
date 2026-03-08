package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/senn404/bookmark-managent/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

// TestPasswordHandler_GenPass verifies that the GenPass handler correctly
// generates and returns a password on success, and returns an appropriate
// error response on failure. It uses a mock Password service.
func TestShortenURLHandler_ShortenURL(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func() *mocks.ShortenURLService

		expectedStatus int
		expectedResp   shortenURLResponse
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {

				body, _ := json.Marshal(shortenURLRequest{
					ExpTime: 100,
					URL:     "huanops.com",
				})

				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", bytes.NewBuffer(body))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func() *mocks.ShortenURLService {
				svcMock := mocks.NewShortenURLService(t)
				svcMock.On("ShortenURL", mock.Anything, mock.Anything, mock.Anything).Return("test", nil)
				return svcMock
			},

			expectedStatus: http.StatusOK,
			expectedResp: shortenURLResponse{
				Code:    "test",
				Message: "Shorten URL generated successfully!",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc()
			testHandler := NewShortenURLHandler(mockSvc)

			testHandler.ShortenURL(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Unmarshal và chuyển về cùng struct để so sánh
			var respon shortenURLResponse
			err := json.Unmarshal(rec.Body.Bytes(), &respon)
			assert.Equal(t, nil, err)
			assert.Equal(t, tc.expectedResp, respon)
		})
	}
}
