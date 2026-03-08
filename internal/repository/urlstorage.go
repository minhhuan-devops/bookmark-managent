package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	urlExpTime = 1 * time.Hour
)

type URLStorage interface {
	StoreURL(ctx context.Context, code, url string, expTime time.Duration) (string, error)
}

type urlStorage struct {
	redisClient *redis.Client
}

func NewURLStorage(redisClent *redis.Client) URLStorage {
	return &urlStorage{
		redisClient: redisClent,
	}
}

//go:generate mockery --name URLStorage --filename URLStorage.go
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, expTime time.Duration) (string, error) {
	return s.redisClient.SetArgs(ctx, code, url, redis.SetArgs{
		Mode: "NX",
		TTL:  expTime,
	}).Result()
}
