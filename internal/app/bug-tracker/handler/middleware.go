package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
)

const (
	authorizationHeader = "Authorization"
	userDataCtx         = "userData"
)

func Logger(logger log.Log) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Internal().
				Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	})
}

func (h *Handler) isUnauthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(authorizationHeader)
		_, err := c.Cookie("refreshToken")

		if header != "" || err == nil {
			return c.JSON(http.StatusBadRequest, newErrorMessage(errUserIsAuthorized))
		}
		return next(c)
	}
}

func (h *Handler) isAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(authorizationHeader)

		if header == "" {
			return c.JSON(http.StatusUnauthorized, newErrorMessage(errInvalidAuthHeader))
		}

		headerParts := strings.Split(header, " ")

		if headerParts[0] != "Bearer" || len(headerParts) != 2 {
			return c.JSON(http.StatusUnauthorized, newErrorMessage(errInvalidAuthHeader))
		}

		if len(headerParts[1]) == 0 {
			return c.JSON(http.StatusUnauthorized, newErrorMessage(errTokenIsEmpty))
		}

		tokenData, err := h.service.Auth.ParseAccessToken(headerParts[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, newErrorMessage(err))
		}

		c.Set(userDataCtx, tokenData)
		return next(c)
	}
}

func getUserData(c echo.Context) (*services.TokenData, error) {
	userData := c.Get(userDataCtx)

	if userData == nil {
		return nil, errUserNotFound
	}

	tokenData, ok := userData.(*services.TokenData)

	if !ok {
		return nil, errUserDataInvalidType
	}

	return tokenData, nil
}
