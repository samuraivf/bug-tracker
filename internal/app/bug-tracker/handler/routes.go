package handler

import (
	"github.com/labstack/echo/v4"
)

func setRoutes(e *echo.Echo, h *Handler) *echo.Echo {
	auth := e.Group(auth)
	{
		auth.POST(signUp, h.signUp)
		auth.POST(signIn, h.signIn, h.isUnauthorized)
		auth.GET(refresh, h.refresh)
		auth.GET(logout, h.logout)
		auth.POST(verify, h.verifyEmail)
		auth.POST(setEmail, h.setEmail)
	}

	return e
}
