package handler

import (
	"fmt"
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

func (h *Handler) refresh(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil || refreshToken == nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidRefreshToken.Error()))
	}

	refreshTokenData, err := h.service.Auth.ParseRefreshToken(refreshToken.Value)

	if err != nil {
		c.SetCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	key := fmt.Sprintf("%s:%s", refreshTokenData.Username, refreshTokenData.TokenID)
	_, err = h.service.Redis.GetRefreshToken(c.Request().Context(), key)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, newErrorMessage(errTokenDoesNotExist.Error()))
	}

	return h.createTokens(c, refreshTokenData.Username, refreshTokenData.UserID)
}

func (h *Handler) logout(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil || refreshToken == nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidRefreshToken.Error()))
	}

	refreshTokenData, err := h.service.Auth.ParseRefreshToken(refreshToken.Value)

	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidRefreshToken.Error()))
	}

	key := fmt.Sprintf("%s:%s", refreshTokenData.Username, refreshTokenData.TokenID)
	err = h.service.Redis.DeleteRefreshToken(c.Request().Context(), key)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err.Error()))
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	
	return c.JSON(http.StatusOK, nil)
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

	key := fmt.Sprintf("%s:%s", username, refreshTokenData.ID)
	err = h.service.Redis.SetRefreshToken(c.Request().Context(), key, refreshTokenData.RefreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err.Error()))
	}

	h.setRefreshToken(c, refreshTokenData.RefreshToken)

	return c.JSON(http.StatusOK, map[string]string{
		"accessToken": accessToken,
	})
}

func (h *Handler) setRefreshToken(c echo.Context, val string) {
	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    val,
		Expires:  time.Now().Add(h.service.Auth.GetRefreshTokenTTL()),
		HttpOnly: true,
	})
}
