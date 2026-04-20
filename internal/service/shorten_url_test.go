package service

import (
	"errors"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/redis/go-redis/v9"
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
			// Case 1: storeurl tra ve ok
			name:    "success",
			url:     "https://huanops.com",
			expTime: 10 * time.Minute,

			setupMock: func() *mocks.URLStorage {
				mockURLStorage := mocks.NewURLStorage(t)
				mockURLStorage.On("StoreURL",
					mock.Anything,
					mock.Anything,
					"https://huanops.com",
					10*time.Minute,
				).Return("OK", nil)
				return mockURLStorage
			},

			expectedErr: nil,
			expectedLen: urlLength, // 9 ký tự
		},
		{
			// Case 2: storeurl tra ve error
			name:    "store url returns error",
			url:     "https://huanops.com",
			expTime: 10 * time.Minute,

			setupMock: func() *mocks.URLStorage {
				mockURLStorage := mocks.NewURLStorage(t)
				mockURLStorage.On("StoreURL",
					mock.Anything,
					mock.Anything,
					"https://huanops.com",
					10*time.Minute,
				).Return("", errors.New("redis connection refused"))
				return mockURLStorage
			},

			expectedErr: errors.New("redis connection refused"),
			expectedLen: 0,
		},
		{
			// Case 3: storeurl tra ve "" lan dau (code bi trung, NX fail),
			// roi tra ve "OK" lan thu 2 -> retry thanh cong
			name:    "retry on code collision",
			url:     "https://huanops.com",
			expTime: 5 * time.Minute,

			setupMock: func() *mocks.URLStorage {
				mockURLStorage := mocks.NewURLStorage(t)
				mockURLStorage.On("StoreURL",
					mock.Anything,
					mock.Anything,
					"https://huanops.com",
					5*time.Minute,
				).Return("", nil).Once() // Lần 1: code bị trùng → trả về ""

				mockURLStorage.On("StoreURL",
					mock.Anything,
					mock.Anything,
					"https://huanops.com",
					5*time.Minute,
				).Return("OK", nil).Once() // Lần 2: code mới → trả về "OK"

				return mockURLStorage
			},

			expectedErr: nil,
			expectedLen: urlLength,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockURLStorage := tc.setupMock()
			testSvc := NewShortenURLService(mockURLStorage)

			respon, err := testSvc.ShortenURL(t.Context(), tc.url, tc.expTime)

			assert.Equal(t, tc.expectedLen, len(respon))

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		inputCode string

		setupMock func() *mocks.URLStorage

		expectedErr error
		expectedURL string
	}{
		{
			// Case 1: code ton tai, tra ve link
			name:      "success",
			inputCode: "aB3xK9mNp",

			setupMock: func() *mocks.URLStorage {
				mockURLStorage := mocks.NewURLStorage(t)
				mockURLStorage.On("GetURL",
					mock.Anything, // ctx
					"aB3xK9mNp",   // code
				).Return("https://huanops.com", nil)
				return mockURLStorage
			},

			expectedErr: nil,
			expectedURL: "https://huanops.com",
		},
		{
			// Case 2: code khong ton tai
			name:      "code not exists",
			inputCode: "notfound",

			setupMock: func() *mocks.URLStorage {
				mockURLStorage := mocks.NewURLStorage(t)
				mockURLStorage.On("GetURL",
					mock.Anything,
					"notfound",
				).Return("", redis.Nil)
				return mockURLStorage
			},

			expectedErr: ErrCodeNotExist,
			expectedURL: "",
		},
		{
			// Case 3: redis tra ve loi ket noi
			name:      "redis unexpected error",
			inputCode: "abc123xyz",

			setupMock: func() *mocks.URLStorage {
				mockURLStorage := mocks.NewURLStorage(t)
				mockURLStorage.On("GetURL",
					mock.Anything,
					"abc123xyz",
				).Return("", errors.New("redis connection timeout"))
				return mockURLStorage
			},

			expectedErr: errors.New("redis connection timeout"),
			expectedURL: "",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockURLStorage := tc.setupMock()
			testSvc := NewShortenURLService(mockURLStorage)

			respon, err := testSvc.GetURL(t.Context(), tc.inputCode)

			assert.Equal(t, tc.expectedURL, respon)

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}
