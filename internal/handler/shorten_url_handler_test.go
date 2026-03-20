package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/senn404/bookmark-managent/internal/service"
	"github.com/senn404/bookmark-managent/internal/service/mocks"
)

// TestPasswordHandler_GenPass verifies that the GenPass handler correctly
// generates and returns a password on success, and returns an appropriate
// error response on failure. It uses a mock Password service.
func TestShortenURLHandler_ShortenURL(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mocks.ShortenURLService

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
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURLService {
				svcMock := mocks.NewShortenURLService(t)
				svcMock.On("ShortenURL", ctx, "huanops.com", time.Duration(100)*time.Second).Return("test", nil)
				return svcMock
			},

			expectedStatus: http.StatusOK,
			expectedResp: shortenURLResponse{
				Code:    "test",
				Message: "Shorten URL generated successfully!",
			},
		},
		{
			name: "service return error",

			setupRequest: func(ctx *gin.Context) {

				body, _ := json.Marshal(shortenURLRequest{
					ExpTime: 100,
					URL:     "huanops.com",
				})

				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", bytes.NewBuffer(body))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURLService {
				svcMock := mocks.NewShortenURLService(t)
				svcMock.On("ShortenURL", ctx, "huanops.com", time.Duration(100)*time.Second).Return("", errors.New("error"))
				return svcMock
			},

			expectedStatus: http.StatusInternalServerError,
			expectedResp: shortenURLResponse{
				Code:    "",
				Message: "internal server error",
			},
		},
		{
			name: "wrong input",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURLService {
				svcMock := mocks.NewShortenURLService(t)
				return svcMock
			},

			expectedStatus: http.StatusBadRequest,
			expectedResp: shortenURLResponse{
				Code:    "",
				Message: "invaild input",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(gc)
			testHandler := NewShortenURLHandler(mockSvc)

			testHandler.ShortenURL(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusOK {
				// Unmarshal và chuyển về cùng struct để so sánh
				var respon shortenURLResponse
				err := json.Unmarshal(rec.Body.Bytes(), &respon)
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedResp.Message, respon.Message)
				assert.Equal(t, tc.expectedResp.Code, respon.Code)
			} else {
				var respon errorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &respon)
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedResp.Message, respon.Error)
			}

		})
	}
}

func TestShortenURLHandler_GetURL(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mocks.ShortenURLService

		expectedStatus int
		expectedURL    string
		expectedError  string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/links/redirect/123456",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "123456"}}

			},
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURLService {
				svcMock := mocks.NewShortenURLService(t)
				svcMock.On("GetURL", ctx, "123456").Return("https://huanops.com", nil)
				return svcMock
			},
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    "https://huanops.com",
		},
		{
			name: "code not exists",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/links/redirect/notfound",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "notfound"}}
			},
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURLService {
				svcMock := mocks.NewShortenURLService(t)
				svcMock.On("GetURL", ctx, "notfound").Return("", service.ErrCodeNotExist)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "code not exists",
		},
		{
			name: "service internal error",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/links/redirect/abc123",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "abc123"}}
			},
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURLService {
				svcMock := mocks.NewShortenURLService(t)
				svcMock.On("GetURL", ctx, "abc123").Return("", errors.New("redis connection error"))
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal server error",
		},
		{
			name: "empty code",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/",
					nil,
				)
				// Không set ctx.Params -> c.Param("code") trả về ""
			},
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenURLService {
				svcMock := mocks.NewShortenURLService(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invaild input",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(gc)
			testHandler := NewShortenURLHandler(mockSvc)

			testHandler.GetURL(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusMovedPermanently {
				assert.Equal(t, tc.expectedURL, rec.Header().Get("Location"))
			} else {
				var respon errorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &respon)
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedError, respon.Error)
			}
		})
	}
}
