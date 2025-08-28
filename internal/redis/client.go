package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(c *redis.Client) *Client {
	return &Client{
		client: c,
	}
}

type Client struct {
	client *redis.Client
}

func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	_, err := c.client.Set(ctx, key, value, ttl).Result()
	return err
}

func (c *Client) Pop(ctx context.Context, key string) error {
	return c.client.GetDel(ctx, key).Err()
}
