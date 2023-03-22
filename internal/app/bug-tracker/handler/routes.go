package handler

import (
	"context"

	"github.com/labstack/echo/v4"
)

func setRoutes(e *echo.Echo, h *Handler) *echo.Echo {
	auth := e.Group(auth)
	{
		auth.POST(signUp, h.signUp)
		auth.POST(signIn, h.signIn, h.isUnauthorized)
		auth.GET(refresh, h.refresh)
		auth.GET(logout, h.logout)
	}

	e.GET("/hello", func(c echo.Context) error {
		data, err := getUserData(c)
		if err != nil {
			return c.String(404, err.Error())
		}

		message := data.Username
		err = h.kafka.Write(message)
		if err == nil {
			h.log.Infof("Sent message: %s", message)
		} else if err == context.Canceled {
			h.log.Error(err)
			return c.String(500, "context canceled")
		} else {
			h.log.Error(err)
			return c.String(500, err.Error())
		}

		return c.String(200, data.Username)
	}, h.isAuthorized)

	return e
}