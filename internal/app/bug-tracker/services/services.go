package services

import (
	"github.com/rs/zerolog"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
)

type Auth interface {
	GenerateAccessToken(username string, userID uint64) (string, error)
	GenerateRefreshToken(username string, userID uint64) (string, error)
	ParseAccessToken(accessToken string) (*TokenData, error)
	ParseRefreshToken(refreshToken string) (*TokenData, error)
}

type User interface {
	GetUserByEmail()
	GetUserById()
	CreateUser()
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
