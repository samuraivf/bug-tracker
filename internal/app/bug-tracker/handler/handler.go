package handler

import (
	"net/http"
	"time"

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

	if _, err := h.service.User.GetUserByEmail(userData.Email); err == nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errUserEmailAlreadyExists.Error()))
	}

	if _, err := h.service.User.GetUserByUsername(userData.Email); err == nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errUserUsernameAlreadyExists.Error()))
	}

	id, err := h.service.User.CreateUser(userData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(errInternalServerError.Error()))
	}

	return c.JSON(http.StatusOK, id)
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

	user, err := h.service.User.ValidateUser(userData.Email, userData.Password)
	if err != nil {
		h.log.Error().Err(err).Msg("")
		c.JSON(http.StatusBadRequest, newErrorMessage(err.Error()))
	}

	return h.createTokens(c, user.Username, user.ID)
}

func (h *Handler) createTokens(c echo.Context, username string, userID uint64) error {
	accessToken, err := h.service.Auth.GenerateAccessToken(username, userID)
	if err != nil {
		h.log.Error().Err(err).Msg("")
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err.Error()))
	}

	refreshTokenData, err := h.service.Auth.GenerateRefreshToken(username, userID)
	if err != nil {
		h.log.Error().Err(err).Msg("")
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err.Error()))
	}

	c.SetCookie(&http.Cookie{
		Name: "refreshToken",
		Value: refreshTokenData,
		Expires: time.Now().Add(h.service.Auth.GetRefreshTokenTTL()),
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, map[string]string{
		"accessToken": accessToken,
	})
}
