package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
)

type Handler struct {
	service *services.Service
	log     *zerolog.Logger
}

func NewHandler(s *services.Service, log *zerolog.Logger) *Handler {
	return &Handler{s, log}
}

func (h *Handler) signUp(c echo.Context) error {
	userData := new(dto.SignUpDto)

	if err := c.Bind(userData); err != nil {
		h.log.Error().Err(err).Msg("")
		return c.JSON(http.StatusBadRequest, newErrorMessage(err.Error()))
	}

	if err := c.Validate(userData); err != nil {
		h.log.Error().Err(err).Msg("")
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidSignUpData.Error()))
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) signIn(c echo.Context) error {
	userData := new(dto.SignInDto)

	if err := c.Bind(userData); err != nil {
		h.log.Error().Err(err).Msg("")
		return c.JSON(http.StatusBadRequest, newErrorMessage(err.Error()))
	}

	if err := c.Validate(userData); err != nil {
		h.log.Error().Err(err).Msg("")
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidSignInData.Error()))
	}

	return c.JSON(http.StatusOK, nil)
}
