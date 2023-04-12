package handler

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
)

func CreateServer() {
	logger := log.New()
	e := echo.New()
	validator := validator.New()
	e.Validator = newValidator(validator)

	dep, close := createDependencies(logger)
	defer close()

	repo := repository.NewRepository(dep.db, logger)
	s := services.NewService(repo, dep.redis)

	h := NewHandler(s, logger, dep.kafka)

	e.Use(Logger(logger))
	e.Use(middleware.Recover())
	e = setRoutes(e, h)

	go func() {
		if err := e.Start(":" + viper.GetString("server-port")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	logger.Info("Server started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	logger.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	logger.Info("Server Exited Properly")
}
