package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found in cache")
	ErrGetKey      = errors.New("failed to get key from redis cache")
	ErrPingTimeout = errors.New("redis ping timeout")
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, ttl time.Duration) (*RedisCache, error) {
	c := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := c.Ping(ctx)

	if status.Err() != nil {
		return nil, ErrPingTimeout
	}

	return &RedisCache{
		client: c,
		ttl:    ttl,
	}, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, r.ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string, into interface{}) error {
	res, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrKeyNotFound
		}

		return ErrGetKey
	}

	return json.Unmarshal(res, &into)
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()

	return err
}

func (r *RedisCache) Stop() error {
	return r.client.Close()
}
