package services

type Service struct {
	Auth
}

type Auth interface {
	GenerateAccessToken(username string, userID uint64) (string, error)
	GenerateRefreshToken(username string, userID uint64) (string, error)
	ParseAccessToken(accessToken string) (*TokenData, error)
	ParseRefreshToken(refreshToken string) (*TokenData, error)
}

func NewService() *Service {
	return &Service{
		Auth: NewAuth(),
	}
}
