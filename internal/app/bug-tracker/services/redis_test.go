package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_redis "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis/mocks"
	"github.com/stretchr/testify/require"
)

func Test_Set(t *testing.T) {
	err := errors.New("err")

	tests := []struct {
		name          string
		key           string
		val           string
		exp           time.Duration
		expectedError error
	}{
		{
			name:          "Error",
			key:           "key",
			val:           "val",
			exp:           time.Minute,
			expectedError: err,
		},
		{
			name:          "OK",
			key:           "key",
			val:           "val",
			exp:           time.Minute,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			ctx := context.Background()

			mock := mock_redis.NewMockRedis(c)
			mock.EXPECT().Set(ctx, test.key, test.val, test.exp).Return(test.expectedError)

			redis := NewRedis(mock)

			actualError := redis.Set(ctx, test.key, test.val, test.exp)

			require.Equal(t, test.expectedError, actualError)
		})
	}
}

func Test_Get(t *testing.T) {
	err := errors.New("err")

	tests := []struct {
		name           string
		key            string
		expectedError  error
		expectedResult string
	}{
		{
			name:           "Error",
			key:            "key",
			expectedError:  err,
			expectedResult: "",
		},
		{
			name:           "OK",
			key:            "key",
			expectedError:  nil,
			expectedResult: "result",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			ctx := context.Background()

			mock := mock_redis.NewMockRedis(c)
			mock.EXPECT().Get(ctx, test.key).Return(test.expectedResult, test.expectedError)

			redis := NewRedis(mock)

			actualResult, actualError := redis.Get(ctx, test.key)

			require.Equal(t, test.expectedResult, actualResult)
			require.Equal(t, test.expectedError, actualError)
		})
	}
}

func Test_SetRefreshToken(t *testing.T) {
	err := errors.New("err")

	tests := []struct {
		name          string
		key           string
		refreshToken  string
		expectedError error
	}{
		{
			name:          "Error",
			key:           "key",
			refreshToken:  "token",
			expectedError: err,
		},
		{
			name:          "OK",
			key:           "key",
			refreshToken:  "token",
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			ctx := context.Background()

			mock := mock_redis.NewMockRedis(c)
			mock.EXPECT().SetRefreshToken(ctx, test.key, test.refreshToken, refreshTokenTTL).Return(test.expectedError)

			redis := NewRedis(mock)

			actualError := redis.SetRefreshToken(ctx, test.key, test.refreshToken)

			require.Equal(t, test.expectedError, actualError)
		})
	}
}

func Test_GetRefreshToken(t *testing.T) {
	err := errors.New("err")

	tests := []struct {
		name           string
		key            string
		expectedError  error
		expectedResult string
	}{
		{
			name:           "Error",
			key:            "key",
			expectedError:  err,
			expectedResult: "",
		},
		{
			name:           "OK",
			key:            "key",
			expectedError:  nil,
			expectedResult: "token",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			ctx := context.Background()

			mock := mock_redis.NewMockRedis(c)
			mock.EXPECT().GetRefreshToken(ctx, test.key).Return(test.expectedResult, test.expectedError)

			redis := NewRedis(mock)

			actualResult, actualError := redis.GetRefreshToken(ctx, test.key)

			require.Equal(t, test.expectedResult, actualResult)
			require.Equal(t, test.expectedError, actualError)
		})
	}
}

func Test_DeleteRefreshToken(t *testing.T) {
	err := errors.New("err")

	tests := []struct {
		name          string
		key           string
		expectedError error
	}{
		{
			name:          "Error",
			key:           "key",
			expectedError: err,
		},
		{
			name:          "OK",
			key:           "key",
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			ctx := context.Background()

			mock := mock_redis.NewMockRedis(c)
			mock.EXPECT().DeleteRefreshToken(ctx, test.key).Return(test.expectedError)

			redis := NewRedis(mock)

			actualError := redis.DeleteRefreshToken(ctx, test.key)

			require.Equal(t, test.expectedError, actualError)
		})
	}
}

func Test_Close(t *testing.T) {
	err := errors.New("err")

	tests := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "Error",
			expectedError: err,
		},
		{
			name:          "OK",
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mock := mock_redis.NewMockRedis(c)
			mock.EXPECT().Close().Return(test.expectedError)

			redis := NewRedis(mock)

			actualError := redis.Close()

			require.Equal(t, test.expectedError, actualError)
		})
	}
}
