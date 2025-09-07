package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewClient initializes and returns a new Client instance using the provided redis.Client.
func NewClient(c *redis.Client) *Client {
	return &Client{
		client: c,
	}
}

// Client wraps a Redis client to provide higher-level operations for interacting with Redis.
type Client struct {
	client *redis.Client
}

// Set stores a key-value pair in the Redis database with the specified time-to-live (ttl) duration.
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	_, err := c.client.Set(ctx, key, value, ttl).Result()
	return err
}

// Pop removes the value stored at the specified key in the Redis database and returns an error if the operation fails.
func (c *Client) Pop(ctx context.Context, key string) error {
	return c.client.GetDel(ctx, key).Err()
}
