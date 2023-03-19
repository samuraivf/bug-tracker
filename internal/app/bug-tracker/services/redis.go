package services

import (
	"context"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis"
)

type RedisService struct {
	repo redis.Redis
}

func NewRedis(repo redis.Redis) *RedisService {
	return &RedisService{repo}
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
