package handler

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	redisgo "github.com/redis/go-redis/v9"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/kafka"
	mock_kafka "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/kafka/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis"
	mock_redis "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
)

func Test_createDependencies(t *testing.T) {
	type mockBehaviour func(*gomock.Controller) log.Log
	err := errors.New("error")

	var (
		db, dbMock, _ = sqlmock.New()
		redisDB, _    = redismock.NewClientMock()
	)

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
	}{
		{
			name: "Error in openPostgres",
			mockBehaviour: func(c *gomock.Controller) log.Log {
				logger := mock_log.NewMockLog(c)

				openPostgres = func(cfg *repository.PostgresConfig) (*sql.DB, error) {
					return nil, err
				}

				logger.EXPECT().Fatal(err)

				return logger
			},
		},
		{
			name: "Error in newRedisClient",
			mockBehaviour: func(c *gomock.Controller) log.Log {
				logger := mock_log.NewMockLog(c)

				openPostgres = func(cfg *repository.PostgresConfig) (*sql.DB, error) {
					return db, nil
				}

				logger.EXPECT().Info("Open PostgreSQL db connection").Return()

				newRedisClient = func(ctx context.Context, cfg *redis.Config) (*redisgo.Client, error) {
					return nil, err
				}

				logger.EXPECT().Fatal(err)

				return logger
			},
		},
		{
			name: "Error in close func",
			mockBehaviour: func(c *gomock.Controller) log.Log {
				logger := mock_log.NewMockLog(c)

				openPostgres = func(cfg *repository.PostgresConfig) (*sql.DB, error) {
					return db, nil
				}

				logger.EXPECT().Info("Open PostgreSQL db connection").Return()

				newRedisClient = func(ctx context.Context, cfg *redis.Config) (*redisgo.Client, error) {
					return redisDB, nil
				}

				logger.EXPECT().Info("Redis started").Return()

				kafkaMock := mock_kafka.NewMockKafka(c)
				newKafkaWriter = func(cfg kafkago.WriterConfig, log log.Log) kafka.Kafka {
					return kafkaMock
				}

				redisMock := mock_redis.NewMockRedis(c)
				newRedis = func(client *redisgo.Client, log log.Log) redis.Redis {
					return redisMock
				}

				logger.EXPECT().Info("Kafka started").Return()

				dbMock.ExpectClose().WillReturnError(err)
				logger.EXPECT().Error(err).Return()
				redisMock.EXPECT().Close().Return(err)
				logger.EXPECT().Error(err).Return()
				kafkaMock.EXPECT().Close().Return(err)
				logger.EXPECT().Error(err).Return()

				logger.EXPECT().Info("PostgreSQL connection closed").Return()
				logger.EXPECT().Info("Redis connection closed").Return()
				logger.EXPECT().Info("Kafka connection closed").Return()

				return logger
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller) log.Log {
				logger := mock_log.NewMockLog(c)

				openPostgres = func(cfg *repository.PostgresConfig) (*sql.DB, error) {
					return db, nil
				}

				logger.EXPECT().Info("Open PostgreSQL db connection").Return()

				newRedisClient = func(ctx context.Context, cfg *redis.Config) (*redisgo.Client, error) {
					return redisDB, nil
				}

				logger.EXPECT().Info("Redis started").Return()

				redisMock := mock_redis.NewMockRedis(c)
				newRedis = func(client *redisgo.Client, log log.Log) redis.Redis {
					return redisMock
				}

				kafkaMock := mock_kafka.NewMockKafka(c)
				newKafkaWriter = func(cfg kafkago.WriterConfig, log log.Log) kafka.Kafka {
					return kafkaMock
				}

				logger.EXPECT().Info("Kafka started").Return()

				dbMock.ExpectClose()
				redisMock.EXPECT().Close().Return(nil)
				kafkaMock.EXPECT().Close().Return(nil)

				logger.EXPECT().Info("PostgreSQL connection closed").Return()
				logger.EXPECT().Info("Redis connection closed").Return()
				logger.EXPECT().Info("Kafka connection closed").Return()

				return logger
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			oldOpenPostgres := openPostgres
			oldNewRedisClient := newRedisClient
			oldNewKafkaWriter := newKafkaWriter
			oldNewRedis := newRedis
			defer func() {
				openPostgres = oldOpenPostgres
				newRedisClient = oldNewRedisClient
				newRedis = oldNewRedis
				newKafkaWriter = oldNewKafkaWriter
			}()

			logger := test.mockBehaviour(c)

			dp, close := createDependencies(logger)
			close()

			if test.name == "OK" {
				assert.NotNil(t, dp.db)
				assert.NotNil(t, dp.kafka)
				assert.NotNil(t, dp.redis)
			}
		})
	}
}
