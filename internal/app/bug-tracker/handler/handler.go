package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
)

type Handler struct {
	service *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) signUp(c echo.Context) error {
	userData := new(dto.SignUpDto)

	if err := c.Bind(userData); err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err.Error()))
	}

	if err := c.Validate(userData); err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidSignUpData.Error()))
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) signIn(c echo.Context) error {
	userData := new(dto.SignInDto)

	if err := c.Bind(userData); err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err.Error()))
	}

	if err := c.Validate(userData); err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidSignInData.Error()))
	}

	return c.JSON(http.StatusOK, nil)
}
