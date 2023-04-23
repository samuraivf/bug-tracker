package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
)

func Test_NewRedis(t *testing.T) {
	c := gomock.NewController(t)
	log := mock_log.NewMockLog(c)
	client := &redis.Client{}

	expected := &RedisRepository{
		redis: client,
		log:   log,
	}

	require.Equal(t, expected, NewRedis(client, log))
}

func Test_Set(t *testing.T) {
	type arguments struct {
		key string
		val string
		exp time.Duration
	}
	type mockBehaviour func(c *gomock.Controller, args arguments) *RedisRepository
	err := errors.New("error")

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		expectedError error
		args          arguments
	}{
		{
			name: "Error in redis client",
			mockBehaviour: func(c *gomock.Controller, args arguments) *RedisRepository {
				client, mock := redismock.NewClientMock()
				mock.ExpectSet(args.key, args.val, args.exp).SetErr(err)
				log := mock_log.NewMockLog(c)
				log.EXPECT().Error(err).Return()

				return &RedisRepository{redis: client, log: log}
			},
			expectedError: err,
			args: arguments{
				key: "key",
				val: "val",
				exp: time.Minute,
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, args arguments) *RedisRepository {
				client, mock := redismock.NewClientMock()
				mock.ExpectSet(args.key, args.val, args.exp).SetVal("")
				log := mock_log.NewMockLog(c)
				log.EXPECT().Infof("[Redis] Set %s:%s", args.key, args.val).Return()

				return &RedisRepository{redis: client, log: log}
			},
			expectedError: nil,
			args: arguments{
				key: "key",
				val: "val",
				exp: time.Minute,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			redis := test.mockBehaviour(c, test.args)

			require.Equal(t, test.expectedError, redis.Set(context.Background(), test.args.key, test.args.val, test.args.exp))
		})
	}
}

func Test_Get(t *testing.T) {
	type arguments struct {
		key string
	}
	type mockBehaviour func(c *gomock.Controller, args arguments) *RedisRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		expectedError  error
		expectedResult string
		args           arguments
	}{
		{
			name: "Error in redis client",
			mockBehaviour: func(c *gomock.Controller, args arguments) *RedisRepository {
				client, mock := redismock.NewClientMock()
				mock.ExpectGet(args.key).SetErr(err)
				log := mock_log.NewMockLog(c)
				log.EXPECT().Error(err).Return()

				return &RedisRepository{redis: client, log: log}
			},
			expectedError:  err,
			expectedResult: "",
			args: arguments{
				key: "key",
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, args arguments) *RedisRepository {
				client, mock := redismock.NewClientMock()
				mock.ExpectGet(args.key).SetVal("get result")

				return &RedisRepository{redis: client, log: nil}
			},
			expectedError:  nil,
			expectedResult: "get result",
			args: arguments{
				key: "key",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			redis := test.mockBehaviour(c, test.args)
			val, err := redis.Get(context.Background(), test.args.key)

			require.Equal(t, test.expectedError, err)
			require.Equal(t, test.expectedResult, val)
		})
	}
}

func Test_SetRefreshToken(t *testing.T) {
	type args struct {
		key          string
		refreshToken string
		TTL          time.Duration
	}

	type mockBehaviour func(c *gomock.Controller, args args) *RedisRepository
	err := errors.New("error")

	tests := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name: "Error in getUserRefreshTokens",
			args: args{
				key:          "key:abc",
				refreshToken: "token",
				TTL:          time.Minute,
			},
			mockBehaviour: func(c *gomock.Controller, args args) *RedisRepository {
				db, mock := redismock.NewClientMock()

				pattern := fmt.Sprintf("*%s*", strings.Split(args.key, ":")[0])
				mock.ExpectKeys(pattern).SetErr(err)

				return &RedisRepository{redis: db}
			},
			expectedError: err,
		},
		{
			name: "Error in deleteUserRefreshTokens",
			args: args{
				key:          "key:abc",
				refreshToken: "token",
				TTL:          time.Minute,
			},
			mockBehaviour: func(c *gomock.Controller, args args) *RedisRepository {
				db, mock := redismock.NewClientMock()

				pattern := fmt.Sprintf("*%s*", strings.Split(args.key, ":")[0])
				tokens := []string{"a", "b", "c", "d", "e", "f"}

				mock.ExpectKeys(pattern).SetVal(tokens)
				mock.ExpectDel(tokens...).SetErr(err)

				return &RedisRepository{redis: db}
			},
			expectedError: err,
		},
		{
			name: "Error in Set",
			args: args{
				key:          "key:abc",
				refreshToken: "token",
				TTL:          time.Minute,
			},
			mockBehaviour: func(c *gomock.Controller, args args) *RedisRepository {
				db, mock := redismock.NewClientMock()
				log := mock_log.NewMockLog(c)

				pattern := fmt.Sprintf("*%s*", strings.Split(args.key, ":")[0])
				tokens := []string{"a", "b", "c", "d", "e", "f"}

				mock.ExpectKeys(pattern).SetVal(tokens)
				mock.ExpectDel(tokens...).SetVal(6)
				mock.ExpectSet(args.key, args.refreshToken, args.TTL).SetErr(err)
				log.EXPECT().Error(err).Return()

				return &RedisRepository{redis: db, log: log}
			},
			expectedError: err,
		},
		{
			name: "OK",
			args: args{
				key:          "key:abc",
				refreshToken: "token",
				TTL:          time.Minute,
			},
			mockBehaviour: func(c *gomock.Controller, args args) *RedisRepository {
				db, mock := redismock.NewClientMock()
				log := mock_log.NewMockLog(c)

				pattern := fmt.Sprintf("*%s*", strings.Split(args.key, ":")[0])
				tokens := []string{"a", "b"}

				mock.ExpectKeys(pattern).SetVal(tokens)
				mock.ExpectSet(args.key, args.refreshToken, args.TTL).SetVal("")
				log.EXPECT().Infof("Set refresh token. Key: %s", args.key).Return()

				return &RedisRepository{redis: db, log: log}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			redis := test.mockBehaviour(c, test.args)

			err := redis.SetRefreshToken(context.Background(), test.args.key, test.args.refreshToken, test.args.TTL)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_getUserRefreshTokens(t *testing.T) {
	type mockBehaviour func(pattern string, expectedResult []string, expectedError error) *RedisRepository

	tests := []struct {
		name           string
		expectedResult []string
		expectedError  error
		pattern        string
		mockBehaviour  mockBehaviour
	}{
		{
			name:          "Error",
			expectedError: errors.New("error"),
			pattern:       "all",
			mockBehaviour: func(pattern string, expectedResult []string, expectedError error) *RedisRepository {
				db, mock := redismock.NewClientMock()

				mock.ExpectKeys(pattern).SetErr(expectedError)

				return &RedisRepository{redis: db}
			},
		},
		{
			name:           "OK",
			expectedError:  nil,
			expectedResult: []string{"res"},
			pattern:        "all",
			mockBehaviour: func(pattern string, expectedResult []string, expectedError error) *RedisRepository {
				db, mock := redismock.NewClientMock()

				mock.ExpectKeys(pattern).SetVal(expectedResult)

				return &RedisRepository{redis: db}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			redis := test.mockBehaviour(test.pattern, test.expectedResult, test.expectedError)

			val, err := redis.getUserRefreshTokens(context.Background(), "all")

			require.Equal(t, test.expectedResult, val)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_deleteUserRefreshTokens(t *testing.T) {
	type mockBehaviour func(keys []string, expectedError error) *RedisRepository

	tests := []struct {
		name          string
		expectedError error
		keys          []string
		mockBehaviour mockBehaviour
	}{
		{
			name:          "Error",
			expectedError: errors.New("error"),
			keys:          []string{"all"},
			mockBehaviour: func(keys []string, expectedError error) *RedisRepository {
				db, mock := redismock.NewClientMock()

				mock.ExpectDel(keys...).SetErr(expectedError)

				return &RedisRepository{redis: db}
			},
		},
		{
			name:          "OK",
			expectedError: nil,
			keys:          []string{"all"},
			mockBehaviour: func(keys []string, expectedError error) *RedisRepository {
				db, mock := redismock.NewClientMock()

				mock.ExpectDel(keys...).SetVal(1)

				return &RedisRepository{redis: db}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			redis := test.mockBehaviour(test.keys, test.expectedError)

			err := redis.deleteUserRefreshTokens(context.Background(), test.keys)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_GetRefreshToken(t *testing.T) {
	type mockBehaviour func(key string) *RedisRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		key            string
		mockBehaviour  mockBehaviour
		expectedResult string
		expectedError  error
	}{
		{
			name: "Error",
			key:  "key",
			mockBehaviour: func(key string) *RedisRepository {
				db, mock := redismock.NewClientMock()

				mock.ExpectGet(key).SetErr(err)

				return &RedisRepository{redis: db}
			},
			expectedResult: "",
			expectedError:  err,
		},
		{
			name: "OK",
			key:  "key",
			mockBehaviour: func(key string) *RedisRepository {
				db, mock := redismock.NewClientMock()

				mock.ExpectGet(key).SetVal("res")

				return &RedisRepository{redis: db}
			},
			expectedResult: "res",
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			redis := test.mockBehaviour(test.key)

			val, err := redis.GetRefreshToken(context.Background(), test.key)

			require.Equal(t, test.expectedResult, val)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_DeleteRefreshToken(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, key string) *RedisRepository
	err := errors.New("error")

	tests := []struct {
		name          string
		key           string
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name: "Error",
			key:  "key",
			mockBehaviour: func(c *gomock.Controller, key string) *RedisRepository {
				db, mock := redismock.NewClientMock()
				log := mock_log.NewMockLog(c)

				mock.ExpectDel(key).SetErr(err)
				log.EXPECT().Error(err)

				return &RedisRepository{redis: db, log: log}
			},
			expectedError: err,
		},
		{
			name: "OK",
			key:  "key",
			mockBehaviour: func(c *gomock.Controller, key string) *RedisRepository {
				db, mock := redismock.NewClientMock()
				log := mock_log.NewMockLog(c)

				mock.ExpectDel(key).SetVal(1)
				log.EXPECT().Infof("Deleted refresh token. Key: %s", key)

				return &RedisRepository{redis: db, log: log}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			redis := test.mockBehaviour(c, test.key)

			err := redis.DeleteRefreshToken(context.Background(), test.key)

			require.Equal(t, test.expectedError, err)
		})
	}
}
