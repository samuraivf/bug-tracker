package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenTTL  = time.Hour * 24
	RefreshTokenTTL = AccessTokenTTL * 30
	jwtAccessKey    = "asdfdgfhaseh281b"
	jwtRefreshKey   = "sgdhksakjajhsdfh"
)

var (
	errInvalidSigningMethod   = errors.New("error invalid signing method")
	errTokenClaimsInvalidType = errors.New("error token claims are not of type *TokenClaims")
)

type AuthService struct{}

type TokenData struct {
	Username string `json:"username"`
	UserID   uint64 `json:"userId"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	TokenData
}

func NewAuth() *AuthService {
	return &AuthService{}
}

func (s *AuthService) GenerateAccessToken(username string, userID uint64) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenData{
			Username: username,
			UserID:   userID,
		},
	})

	return accessToken.SignedString([]byte(jwtAccessKey))
}

func (s *AuthService) GenerateRefreshToken(username string, userID uint64) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenData{
			Username: username,
			UserID:   userID,
		},
	})

	return refreshToken.SignedString([]byte(jwtRefreshKey))
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
