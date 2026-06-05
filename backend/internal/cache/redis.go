package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/config"
)

// Client is a thin wrapper over go-redis with JSON helpers and pattern delete.
type Client struct {
	rdb *redis.Client
}

func New(cfg config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Client{rdb: rdb}, nil
}

func (c *Client) Raw() *redis.Client { return c.rdb }

func (c *Client) Get(ctx context.Context, key string, dst any) (bool, error) {
	b, err := c.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(b, dst)
}

func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, b, ttl).Err()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.rdb.Del(ctx, keys...).Err()
}

// DelByPattern scans + deletes any matching key. Use sparingly.
func (c *Client) DelByPattern(ctx context.Context, pattern string) error {
	iter := c.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	var batch []string
	for iter.Next(ctx) {
		batch = append(batch, iter.Val())
		if len(batch) >= 500 {
			if err := c.rdb.Del(ctx, batch...).Err(); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(batch) > 0 {
		return c.rdb.Del(ctx, batch...).Err()
	}
	return nil
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Incr(ctx, key).Result()
}

// SAdd / SMembers helpers used by chat presence.
func (c *Client) SAdd(ctx context.Context, key string, members ...string) error {
	args := make([]any, len(members))
	for i, m := range members {
		args[i] = m
	}
	return c.rdb.SAdd(ctx, key, args...).Err()
}

func (c *Client) SRem(ctx context.Context, key string, members ...string) error {
	args := make([]any, len(members))
	for i, m := range members {
		args[i] = m
	}
	return c.rdb.SRem(ctx, key, args...).Err()
}

func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.rdb.SMembers(ctx, key).Result()
}

// Publish for the pub/sub chat channel.
func (c *Client) Publish(ctx context.Context, channel string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.rdb.Publish(ctx, channel, b).Err()
}

func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.rdb.Subscribe(ctx, channels...)
}
