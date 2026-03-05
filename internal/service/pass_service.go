package service

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	// charset defines the set of characters used for password generation,
	// including lowercase, uppercase, digits, and special characters.
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// passLength defines the fixed length of generated passwords.
	passLength = 16
)

// passwordService is the concrete implementation of the Password interface.
// It generates cryptographically secure random passwords.
type passwordService struct {
}

// Password defines the interface for password generation operations.
//
//go:generate mockery --name Password --filename pass_service.go
type Password interface {
	// GeneratePassword creates a new random password and returns it.
	// It returns an error if the cryptographic random number generator fails.
	GeneratePassword() (string, error)
}

// NewPassword creates a new Password service instance.
func NewPassword() Password {
	return &passwordService{}
}

// GeneratePassword generates a cryptographically secure random password
// of length passLength using characters from the defined charset.
// It uses crypto/rand for secure random number generation, making the
// generated passwords suitable for security-sensitive applications.
func (s *passwordService) GeneratePassword() (string, error) {
	var strBuilder bytes.Buffer

	for i := 1; i <= passLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}

	return strBuilder.String(), nil
}
