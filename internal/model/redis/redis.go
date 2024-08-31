package redisservice

import (
	"context"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	redis *redis.Client
}

func New(redis *redis.Client) *RedisService {
	return &RedisService{redis: redis}
}

func (r *RedisService) GetByFunc(ctx context.Context, key string, ttl time.Duration, f func() (string, error)) (string, error) {
	v, err := r.redis.Get(ctx, key).Result()
	if err == nil { // cache hit
		return v, nil
	}
	if err != redis.Nil {
		logger.Error("Failed to get value from cache", "err", err)
	}

	// get value from function
	res, err := f()
	if err != nil {
		return "", err
	}

	if err := r.redis.Set(ctx, key, res, ttl).Err(); err != nil {
		logger.Error("Failed to cache value", "err", err)
	}

	return res, nil
}

func (r *RedisService) GetDel(ctx context.Context, key string) (string, error) {
	v, err := r.redis.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return "", domain.ErrCacheMiss
	}
	if err != nil {
		return "", err
	}

	return v, nil
}

func (r *RedisService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.redis.Set(ctx, key, value, ttl).Err()
}
