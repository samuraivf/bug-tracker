package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getUserById(c echo.Context) error {
	id, err := h.params.GetIdParam(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	user, err := h.service.User.GetUserById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, newErrorMessage(errUserNotFound))
	}

	return c.JSON(http.StatusFound, user)
}

func (h *Handler) getUserByUsername(c echo.Context) error {
	username, err := h.params.GetUsernameParam(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	user, err := h.service.User.GetUserByUsername(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, newErrorMessage(errUserNotFound))
	}

	return c.JSON(http.StatusFound, user)
}

func (h *Handler) getUserProjects(c echo.Context) error {
	userData, err := getUserData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	projects, err := h.service.Project.GetProjectsByUserId(userData.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, newErrorMessage(err))
	}

	return c.JSON(http.StatusFound, projects)
}
