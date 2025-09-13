package jwt_test

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/riabininkf/http-auth-example/internal/jwt"
	"github.com/riabininkf/http-auth-example/internal/jwt/mocks"
)

func TestStorage_Save(t *testing.T) {
	t.Run("failed to save into cache", func(t *testing.T) {
		cache := mocks.NewCache(t)
		cache.On("Set", t.Context(), hashStorageKey("test_key"), "", time.Second*5).Return(assert.AnError)

		storage := jwt.NewStorage(time.Second*5, cache)
		assert.Equal(t, assert.AnError, storage.Save(t.Context(), "test_key"))
	})

	t.Run("positive case", func(t *testing.T) {
		cache := mocks.NewCache(t)
		cache.On("Set", t.Context(), hashStorageKey("test_key"), "", time.Second*5).Return(nil)

		storage := jwt.NewStorage(time.Second*5, cache)
		assert.NoError(t, storage.Save(t.Context(), "test_key"))
	})
}

func TestStorage_Pop(t *testing.T) {
	t.Run("failed to pop from cache", func(t *testing.T) {
		cache := mocks.NewCache(t)
		cache.On("Pop", t.Context(), hashStorageKey("test_key")).Return(assert.AnError)

		storage := jwt.NewStorage(time.Second*5, cache)
		assert.Equal(t, assert.AnError, storage.Pop(t.Context(), "test_key"))
	})

	t.Run("positive case", func(t *testing.T) {
		cache := mocks.NewCache(t)
		cache.On("Pop", t.Context(), hashStorageKey("test_key")).Return(nil)

		storage := jwt.NewStorage(time.Second*5, cache)
		assert.NoError(t, storage.Pop(t.Context(), "test_key"))
	})
}

func hashStorageKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return string(sum[:])
}
