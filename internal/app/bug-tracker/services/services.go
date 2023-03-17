package services

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
)

type Auth interface {
	GetRefreshTokenTTL() time.Duration
	GenerateAccessToken(username string, userID uint64) (string, error)
	GenerateRefreshToken(username string, userID uint64) (string, error)
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

type Service struct {
	Auth
	User
}

func NewService(repo *repository.Repository, log *zerolog.Logger) *Service {
	return &Service{
		Auth: NewAuth(log),
		User: NewUser(repo.User),
	}
}
