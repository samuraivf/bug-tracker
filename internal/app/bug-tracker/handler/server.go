package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
)

func CreateServer() {
	e := echo.New()
	e.Validator = newValidator()

	s := services.NewService()
	h := NewHandler(s)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST(authSignUp, h.signUp)
	e.POST(authSignIn, h.signIn)

	e.Logger.Fatal(e.Start(":8080"))
}
