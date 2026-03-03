package service

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestPasswordService_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		expectedLen int
		expectedErr error
	}{
		{
			name: "normal case",

			expectedLen: 16,
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewPassword()

			pass, err := testSvc.GeneratePassword()

			assert.Equal(t, tc.expectedLen, len(pass))
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
