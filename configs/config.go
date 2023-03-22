package configs

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func init() {
	if err := initConfig(); err != nil {
		log.Fatal().Timestamp().Err(err).Msg("")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal().Timestamp().Err(err).Msg("")
	}
}

func PostgresConfig() *repository.PostgresConfig {
	return &repository.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.name"),
	}
}

func RedisConfig() *redis.Config {
	return &redis.Config{
		Host: viper.GetString("redis.host"),
		Port: viper.GetString("redis.port"),
	}
}

func KafkaConfig() kafkago.WriterConfig {
	return kafkago.WriterConfig{
		Brokers: []string{viper.GetString("kafka.brokers")},
		Topic:   viper.GetString("kafka.topic"),
		BatchTimeout: 1 * time.Millisecond,
	}
}