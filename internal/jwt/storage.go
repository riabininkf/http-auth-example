package jwt

import (
	"context"
	"crypto/sha256"
	"time"
)

func NewStorage(
	refreshTokenTTL time.Duration,
	cache Cache,
) *Storage {
	return &Storage{
		refreshTokenTTL: refreshTokenTTL,
		cache:           cache,
	}
}

type (
	Storage struct {
		refreshTokenTTL time.Duration
		cache           Cache
	}

	Cache interface {
		Set(ctx context.Context, key string, value any, ttl time.Duration) error
		Pop(ctx context.Context, key string) error
	}
)

func (s *Storage) Save(ctx context.Context, token string) error {
	return s.cache.Set(ctx, s.hash(token), "", s.refreshTokenTTL)
}

func (s *Storage) Pop(ctx context.Context, token string) error {
	return s.cache.Pop(ctx, s.hash(token))
}

func (s *Storage) hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return string(sum[:])
}
