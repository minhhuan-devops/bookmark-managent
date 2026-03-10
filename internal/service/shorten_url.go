// Package service provides the business logic layer for the Bookmark Management API.
package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/senn404/bookmark-managent/internal/repository"
)

const (
	// urlLength defines the fixed length of the generated short URL code.
	urlLength = 9
)

// shortenURL is the concrete implementation of ShortenURLService.
// It delegates URL storage to a URLStorage repository.
type shortenURL struct {
	urlStorage repository.URLStorage
}

//go:generate mockery --name ShortenURLService --filename shorten_url.go

// ShortenURLService defines the interface for URL shortening operations.
type ShortenURLService interface {
	// ShortenURL stores the given URL with the specified expiration and returns a unique short code.
	// It retries automatically if the generated code already exists in storage.
	ShortenURL(ctx context.Context, url string, expTime time.Duration) (string, error)
}

// NewShortenURLService creates a new ShortenURLService with the given URLStorage.
// It uses dependency injection to allow easy testing with mock repositories.
func NewShortenURLService(urlStorage repository.URLStorage) ShortenURLService {
	return &shortenURL{
		urlStorage: urlStorage,
	}
}

// generateURL generates a random short URL code of length urlLength.
// It uses crypto/rand to produce a cryptographically secure random code.
// It returns an error if the random number generator fails.
func (s *shortenURL) generateURL() (string, error) {
	var urlShorten bytes.Buffer

	for i := 1; i <= urlLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		urlShorten.WriteByte(charset[randomIndex.Int64()])
	}
	return urlShorten.String(), nil
}

// ShortenURL stores the given URL in the repository with the specified expiration duration.
// It generates a unique random code and retries until a collision-free code is stored.
// It returns the generated short code on success, or an error if storage fails.
func (s *shortenURL) ShortenURL(ctx context.Context, url string, expTime time.Duration) (string, error) {
	for {
		urlResponse, err := s.generateURL()
		if err != nil {
			return "", err
		}
		check, err := s.urlStorage.StoreURL(ctx, urlResponse, url, expTime)
		if err != nil {
			return "", err
		}
		if check == "OK" {
			return urlResponse, nil
		}
	}
}
