package redis

import (
	"api-gateway/internal/configs"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(integrationConfig *configs.IntegrationConfig) *RedisRepository {
	return &RedisRepository{
		client: redis.NewClient(&redis.Options{
			Addr:     integrationConfig.RedisAddr,
			Password: integrationConfig.RedisPassword,
			DB:       0,
		}),
	}
}

func (s *RedisRepository) GetUserId(token string) (uuid.UUID, error) {
	ctx := context.Background()
	res := s.client.Get(ctx, token)

	if res.Err() != nil {
		if errors.Is(res.Err(), redis.Nil) {
			return uuid.Nil, nil
		} else {
			return uuid.Nil, fmt.Errorf("failed to load token data from redis: %w", res.Err())
		}
	}

	id, err := uuid.Parse(res.String())
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
