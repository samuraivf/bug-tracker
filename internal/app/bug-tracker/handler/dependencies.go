package handler

import (
	"context"
	"database/sql"

	"github.com/samuraivf/bug-tracker/configs"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/kafka"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
)

type Dependencies struct {
	db    *sql.DB
	redis redis.Redis
	kafka kafka.Kafka
}

var (
	openPostgres   = repository.OpenPostgres
	newRedisClient = redis.NewClient
	newRedis       = redis.NewRedis
	newKafkaWriter = kafka.NewKafkaWriter
)

func createDependencies(logger log.Log) (*Dependencies, func()) {
	db, err := openPostgres(configs.PostgresConfig())
	if err != nil {
		logger.Fatal(err)
		return nil, func() {}
	}
	logger.Info("Open PostgreSQL db connection")

	redisClient, err := newRedisClient(context.Background(), configs.RedisConfig())
	if err != nil {
		logger.Fatal(err)
		return nil, func() {}
	}
	logger.Info("Redis started")

	redisRepo := newRedis(redisClient, logger)

	k := newKafkaWriter(configs.KafkaConfig(), logger)
	logger.Info("Kafka started")

	return &Dependencies{db: db, redis: redisRepo, kafka: k}, func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
		logger.Info("PostgreSQL connection closed")

		if err := redisRepo.Close(); err != nil {
			logger.Error(err)
		}
		logger.Info("Redis connection closed")

		if err := k.Close(); err != nil {
			logger.Error(err)
		}
		logger.Info("Kafka connection closed")
	}
}
