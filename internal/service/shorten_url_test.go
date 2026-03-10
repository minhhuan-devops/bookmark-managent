package service

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/senn404/bookmark-managent/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func TestShortenURL(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name    string
		url     string
		expTime time.Duration

		setupMock func() *mocks.URLStorage

		expectedErr error
		expectedLen int
	}{
		{
			name:    "normal case",
			url:     "huanops.com",
			expTime: 10,

			setupMock: func() *mocks.URLStorage {
				mockURLStorage := mocks.NewURLStorage(t)
				mockURLStorage.On("StoreURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("OK", nil)
				return mockURLStorage
			},

			expectedErr: nil,
			expectedLen: 9,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockURLStorage := tc.setupMock()
			testSvc := NewShortenURLService(mockURLStorage)

			respon, err := testSvc.ShortenURL(t.Context(), tc.url, tc.expTime)

			assert.Equal(t, tc.expectedLen, len(respon))
			assert.Equal(t, tc.expectedErr, err)

		})
	}
}
