package handler

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func CreateServer() {
	e := echo.New()
	h := NewHandler()

	e.Use(middleware.Logger())
  	e.Use(middleware.Recover())

	e.GET("/hello", h.hello)

	e.Logger.Fatal(e.Start(":8080"))
}