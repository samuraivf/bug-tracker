package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Redis interface {
	SetRefreshToken(ctx context.Context, key, refreshToken string, TTL time.Duration) error
	GetRefreshToken(ctx context.Context, key string) (string, error)
	DeleteRefreshToken(ctx context.Context, key string) error
}

type RedisRepository struct {
	redis *redis.Client
	log   *zerolog.Logger
}

func NewClient(ctx context.Context, cfg *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRedis(client *redis.Client, log *zerolog.Logger) Redis {
	return &RedisRepository{
		redis: client,
		log:   log,
	}
}

func (r *RedisRepository) getUserRefreshTokens(ctx context.Context, pattern string) ([]string, error) {
	return r.redis.Keys(ctx, pattern).Result()
}

func (r *RedisRepository) deleteUserRefreshTokens(ctx context.Context, keys []string) error {
	return r.redis.Del(ctx, keys...).Err()
}

func (r *RedisRepository) SetRefreshToken(ctx context.Context, key, refreshToken string, TTL time.Duration) error {
	userTokens, err := r.getUserRefreshTokens(ctx, fmt.Sprintf("*%s*", strings.Split(key, ":")[0]))

	if err != nil {
		return err
	}

	if len(userTokens) > 5 {
		err := r.deleteUserRefreshTokens(ctx, userTokens)

		if err != nil {
			return err
		}
	}

	err = r.redis.Set(ctx, key, refreshToken, TTL).Err()
	if err != nil {
		r.log.Error().Err(err).Msg("")
		return err
	}
	r.log.Info().Msgf("Set refresh token. Key: %s", key)

	return nil
}

func (r *RedisRepository) GetRefreshToken(ctx context.Context, key string) (string, error) {
	refreshToken, err := r.redis.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (r *RedisRepository) DeleteRefreshToken(ctx context.Context, key string) error {
	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		r.log.Error().Err(err).Msg("")
		return err
	}
	r.log.Info().Msgf("Deleted refresh token. Key: %s", key)

	return nil
}
