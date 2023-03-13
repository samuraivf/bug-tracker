package handler

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func CreateServer() {
	e := echo.New()
	e.Validator = newValidator()
	h := NewHandler()

	e.Use(middleware.Logger())
  	e.Use(middleware.Recover())

	e.POST(authSignUp, h.signUp)
	e.POST(authSignIn, h.signIn)

	e.Logger.Fatal(e.Start(":8080"))
}