package handler

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
)

func CreateServer() {
	logger := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := initConfig(); err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	e := echo.New()
	e.Validator = newValidator()

	config := &repository.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.name"),
	}

	db, err := repository.OpenPostgres(config)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	logger.Info().Msg("Open PostgreSQL db connection")
	defer db.Close()

	repo := repository.NewRepository(db, &logger)
	s := services.NewService(repo, &logger)
	h := NewHandler(s, &logger)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))
	e.Use(middleware.Recover())

	e.POST(authSignUp, h.signUp)
	e.POST(authSignIn, h.signIn)

	e.Logger.Fatal(e.Start(":" + viper.GetString("server-port")))
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
