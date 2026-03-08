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
	urlLength = 9
)

type shortenURL struct {
	urlStorage repository.URLStorage
}

type ShortenURLService interface {
	ShortenURL(ctx context.Context, url string, expTime time.Duration) (string, error)
}

func NewShortenURLService(urlStorage repository.URLStorage) ShortenURLService {
	return &shortenURL{
		urlStorage: urlStorage,
	}
}

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

//go:generate mockery --name ShortenURLService --filename shorten_url.go
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
		urlResponse, err = s.generateURL()
	}
}
