package services

import (
	"context"
	"time"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
)

type Auth interface {
	GetRefreshTokenTTL() time.Duration
	GenerateAccessToken(username string, userID uint64) (string, error)
	GenerateRefreshToken(username string, userID uint64) (*RefreshTokenData, error)
	ParseAccessToken(accessToken string) (*TokenData, error)
	ParseRefreshToken(refreshToken string) (*TokenData, error)
}

type User interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id uint64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(userData *dto.SignUpDto) (uint64, error)
	ValidateUser(email, password string) (*models.User, error)
}

type Redis interface {
	SetRefreshToken(ctx context.Context, key, refreshToken string) error
	GetRefreshToken(ctx context.Context, key string) (string, error)
	DeleteRefreshToken(ctx context.Context, key string) error
}

type Service struct {
	Auth
	User
	Redis
}

func NewService(repo *repository.Repository, redisRepo redis.Redis) *Service {
	return &Service{
		Auth:  NewAuth(),
		User:  NewUser(repo.User),
		Redis: NewRedis(redisRepo),
	}
}
