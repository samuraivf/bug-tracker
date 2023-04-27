package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func Test_GetRefreshTokenTTL(t *testing.T) {
	ttl := NewAuth().GetRefreshTokenTTL()

	require.Equal(t, refreshTokenTTL, ttl)
	require.Equal(t, time.Hour*24*30, ttl)
}

func Test_GenerateAccessToken(t *testing.T) {
	auth := NewAuth()
	username := "username"
	userID := uint64(1)

	expectedToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenData{
			Username: username,
			UserID:   userID,
		},
	}).SignedString([]byte(jwtAccessKey))

	token, err := auth.GenerateAccessToken(username, userID)

	require.NoError(t, err)
	require.Equal(t, expectedToken, token)
}

func Test_GenerateRefreshToken(t *testing.T) {
	auth := NewAuth()
	username := "username"
	userID := uint64(1)

	token, err := auth.GenerateRefreshToken(username, userID)

	require.NoError(t, err)
	require.NotEmpty(t, token.ID)

	tokenData, err := auth.ParseRefreshToken(token.RefreshToken)

	require.NoError(t, err)
	require.Equal(t, username, tokenData.Username)
	require.Equal(t, userID, tokenData.UserID)
}

func Test_ParseAccessToken(t *testing.T) {
	auth := NewAuth()
	username := "username"
	userID := uint64(1)

	token, _ := auth.GenerateAccessToken(username, userID)
	tokenData, err := auth.ParseAccessToken(token)

	require.NoError(t, err)
	require.Equal(t, username, tokenData.Username)
	require.Equal(t, userID, tokenData.UserID)
}

func Test_ParseRefreshToken(t *testing.T) {
	auth := NewAuth()
	username := "username"
	userID := uint64(1)

	token, _ := auth.GenerateRefreshToken(username, userID)
	tokenData, err := auth.ParseRefreshToken(token.RefreshToken)

	require.NotEmpty(t, token.ID)
	require.NoError(t, err)
	require.Equal(t, username, tokenData.Username)
	require.Equal(t, userID, tokenData.UserID)
}
