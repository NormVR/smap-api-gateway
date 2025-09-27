package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository() *RedisRepository {
	return &RedisRepository{
		client: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "redispass",
			DB:       0,
		}),
	}
}

func (s *RedisRepository) GetUserId(token string) (int64, error) {
	ctx := context.Background()
	userId, err := s.client.Get(ctx, token).Int64()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		} else {
			return 0, fmt.Errorf("failed to load token data from redis: %w", err)
		}
	}

	return userId, nil
}
