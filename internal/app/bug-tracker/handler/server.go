package handler

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
)

func CreateServer() {
	e := echo.New()
	e.Validator = newValidator()

	logger := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr})

	s := services.NewService(&logger)
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

	e.Logger.Fatal(e.Start(":8080"))
}
