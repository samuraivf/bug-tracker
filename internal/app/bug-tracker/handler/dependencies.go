package handler

import (
	"context"
	"database/sql"

	redisgo "github.com/redis/go-redis/v9"

	"github.com/samuraivf/bug-tracker/configs"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/kafka"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
)

type Dependencies struct {
	db    *sql.DB
	redis *redisgo.Client
	kafka *kafka.KafkaWriter
}

func CreateDependencies(logger log.Log) (*Dependencies, func()) {
	db, err := repository.OpenPostgres(configs.PostgresConfig())
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Open PostgreSQL db connection")

	redisClient, err := redis.NewClient(context.Background(), configs.RedisConfig())
	if err != nil {
		logger.Fatal(err)
	}
	defer redisClient.Conn().Close()

	k := kafka.NewKafkaWriter(configs.KafkaConfig(), logger)
	logger.Info("Kafka started")

	return &Dependencies{db: db, redis: redisClient, kafka: k}, func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
		logger.Info("PostgreSQL connection closed")

		if err := redisClient.Conn().Close(); err != nil {
			logger.Error(err)
		}
		logger.Info("Redis connection closed")

		if err := k.Close(); err != nil {
			logger.Error(err)
		}
		logger.Info("Kafka connection closed")
	}
}
