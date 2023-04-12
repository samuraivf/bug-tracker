package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
)

//go:generate mockgen -source=redis.go -destination=mocks/redis.go

type Redis interface {
	Set(ctx context.Context, key, val string, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	SetRefreshToken(ctx context.Context, key, refreshToken string, TTL time.Duration) error
	GetRefreshToken(ctx context.Context, key string) (string, error)
	DeleteRefreshToken(ctx context.Context, key string) error
	Close() error
}

type RedisRepository struct {
	redis *redis.Client
	log   log.Log
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

func NewRedis(client *redis.Client, log log.Log) Redis {
	return &RedisRepository{
		redis: client,
		log:   log,
	}
}

func (r *RedisRepository) Set(ctx context.Context, key, val string, exp time.Duration) error {
	err := r.redis.Set(ctx, key, val, exp).Err()
	if err != nil {
		r.log.Error(err)
		return err
	}
	r.log.Infof("[Redis] Set %s:%s", key, val)

	return nil
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		r.log.Error(err)
		return "", err
	}
	return val, nil
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
		r.log.Error(err)
		return err
	}
	r.log.Infof("Set refresh token. Key: %s", key)

	return nil
}

func (r *RedisRepository) getUserRefreshTokens(ctx context.Context, pattern string) ([]string, error) {
	return r.redis.Keys(ctx, pattern).Result()
}

func (r *RedisRepository) deleteUserRefreshTokens(ctx context.Context, keys []string) error {
	return r.redis.Del(ctx, keys...).Err()
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
		r.log.Error(err)
		return err
	}
	r.log.Infof("Deleted refresh token. Key: %s", key)

	return nil
}

func (r *RedisRepository) Close() error {
	return r.redis.Conn().Close()
}
