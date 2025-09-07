package jwt

//go:generate mockery --name Cache --output ./mocks --outpkg mocks --filename cache.go --structname Cache

import (
	"context"
	"crypto/sha256"
	"time"
)

// NewStorage initializes a new Storage instance with the provided refresh token TTL and cache implementation.
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
	// Storage represents a storage mechanism with token TTL and an associated cache implementation.
	Storage struct {
		refreshTokenTTL time.Duration
		cache           Cache
	}

	// Cache defines methods for managing a key-value store with optional context and TTL (time-to-live) functionality.
	// Set stores a value for a key in the cache with a specified TTL, returning an error if the operation fails.
	// Pop removes a key and its associated value from the cache, returning an error if the operation fails.
	Cache interface {
		Set(ctx context.Context, key string, value any, ttl time.Duration) error
		Pop(ctx context.Context, key string) error
	}
)

// Save stores the given token in the cache with a configured TTL. Returns an error if the operation fails.
func (s *Storage) Save(ctx context.Context, token string) error {
	return s.cache.Set(ctx, s.hash(token), "", s.refreshTokenTTL)
}

// Pop removes the specified token from the cache using its hashed value and the provided context.
// It returns an error if the operation fails.
func (s *Storage) Pop(ctx context.Context, token string) error {
	return s.cache.Pop(ctx, s.hash(token))
}

// hash generates a SHA-256 hash of the provided token and returns it as a string.
func (s *Storage) hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return string(sum[:])
}
