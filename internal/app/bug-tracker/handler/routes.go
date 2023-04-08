package handler

import (
	"github.com/labstack/echo/v4"
)

func setRoutes(e *echo.Echo, h *Handler) *echo.Echo {
	auth := e.Group(auth)
	{
		auth.POST(signUp, h.signUp)
		auth.POST(signIn, func(c echo.Context) error {
			return h.signIn(c, h.createTokens)
		}, h.isUnauthorized)
		auth.GET(refresh, func(c echo.Context) error {
			return h.refresh(c, h.createTokens)
		})
		auth.GET(logout, h.logout)
		auth.POST(verify, h.verifyEmail)
		auth.POST(setEmail, h.setEmail)
	}

	return e
}
