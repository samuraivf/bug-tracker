package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	accessTokenTTL  = time.Hour * 24
	refreshTokenTTL = accessTokenTTL * 30
	jwtAccessKey    = "asdfdgfhaseh281b"
	jwtRefreshKey   = "sgdhksakjajhsdfh"
)

var (
	errInvalidSigningMethod   = errors.New("error invalid signing method")
	errTokenClaimsInvalidType = errors.New("error token claims are not of type *TokenClaims")
)

type AuthService struct {
	log *zerolog.Logger
}

type TokenData struct {
	TokenID  string `json:"tokenId"`
	Username string `json:"username"`
	UserID   uint64 `json:"userId"`
}

type RefreshTokenData struct {
	ID           string
	RefreshToken string
}

type TokenClaims struct {
	jwt.RegisteredClaims
	TokenData
}

func NewAuth(log *zerolog.Logger) *AuthService {
	return &AuthService{log}
}

func (s *AuthService) GetRefreshTokenTTL() time.Duration {
	return refreshTokenTTL
}

func (s *AuthService) GenerateAccessToken(username string, userID uint64) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenData{
			Username: username,
			UserID:   userID,
		},
	})

	return accessToken.SignedString([]byte(jwtAccessKey))
}

func (s *AuthService) GenerateRefreshToken(username string, userID uint64) (*RefreshTokenData, error) {
	tokenID := uuid.NewString()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenData{
			TokenID: tokenID,
			Username: username,
			UserID:   userID,
		},
	})

	token, err := refreshToken.SignedString([]byte(jwtRefreshKey))
	if err != nil {
		return nil, err
	}

	return &RefreshTokenData{RefreshToken: token, ID: tokenID}, nil
}

func (s *AuthService) ParseAccessToken(accessToken string) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidSigningMethod
		}

		return []byte(jwtAccessKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errTokenClaimsInvalidType
	}

	return &claims.TokenData, nil
}

func (s *AuthService) ParseRefreshToken(refreshToken string) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidSigningMethod
		}

		return []byte(jwtRefreshKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errTokenClaimsInvalidType
	}

	return &claims.TokenData, nil
}
