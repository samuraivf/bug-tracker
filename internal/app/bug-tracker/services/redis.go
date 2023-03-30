package services

import (
	"context"
	"time"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis"
)

type RedisService struct {
	repo redis.Redis
}

func NewRedis(repo redis.Redis) Redis {
	return &RedisService{repo}
}

func (s *RedisService) Set(ctx context.Context, key, val string, exp time.Duration) error {
	return s.repo.Set(ctx, key, val, exp)
}

func (s *RedisService) Get(ctx context.Context, key string) (string, error) {
	return s.repo.Get(ctx, key)
}

func (s *RedisService) SetRefreshToken(ctx context.Context, key, refreshToken string) error {
	return s.repo.SetRefreshToken(ctx, key, refreshToken, refreshTokenTTL)
}

func (s *RedisService) GetRefreshToken(ctx context.Context, key string) (string, error) {
	return s.repo.GetRefreshToken(ctx, key)
}

func (s *RedisService) DeleteRefreshToken(ctx context.Context, key string) error {
	return s.repo.DeleteRefreshToken(ctx, key)
}
