package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	client *redis.Client
}

func NewTokenRepository(redisURL string) (*repository, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &repository{client: client}, nil
}

func (r *repository) StoreToken(ctx context.Context, userID uuid.UUID, token string) error {
	return r.client.Set(ctx, token, userID.String(), 24*time.Hour).Err()
}

func (r *repository) ValidateToken(ctx context.Context, token string) (uuid.UUID, error) {
	userIDStr, err := r.client.Get(ctx, token).Result()
	if err == redis.Nil {
		return uuid.Nil, fmt.Errorf("token not found")
	}
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %v", err)
	}

	return userID, nil
}
