package service

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	passLength = 16
)

type passwordService struct {
}

//go:generate mockery --name Password --filename pass_service.go
type Password interface {
	GeneratePassword() (string, error)
}

func NewPassword() Password {
	return &passwordService{}
}

// @Summary GeneratePassword
// @Description GeneratePassword
// @Tags password
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 500 {string} string
// @Router /gen-pass [get]
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
