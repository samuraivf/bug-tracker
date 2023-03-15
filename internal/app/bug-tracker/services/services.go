package services

import "github.com/rs/zerolog"

type Service struct {
	Auth
}

type Auth interface {
	GenerateAccessToken(username string, userID uint64) (string, error)
	GenerateRefreshToken(username string, userID uint64) (string, error)
	ParseAccessToken(accessToken string) (*TokenData, error)
	ParseRefreshToken(refreshToken string) (*TokenData, error)
}

func NewService(log *zerolog.Logger) *Service {
	return &Service{
		Auth: NewAuth(log),
	}
}
