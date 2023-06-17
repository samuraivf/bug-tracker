package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
)

func (h *Handler) createTask(c echo.Context) error {
	taskData := new(dto.CreateTaskDto)
	userData, err := getUserData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	if err := c.Bind(taskData); err != nil {
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidJSON))
	}

	if err := c.Validate(taskData); err != nil {
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidTaskData))
	}

	id, err := h.service.Task.CreateTask(taskData, userData.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err))
	}

	return c.JSON(http.StatusOK, id)
}

func (h *Handler) workOnTask(c echo.Context) error {
	workOnTaskData := new(dto.WorkOnTaskDto)
	userData, err := getUserData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	if err := c.Bind(workOnTaskData); err != nil {
		return c.JSON(http.StatusBadRequest, errInvalidJSON)
	}

	if err := c.Validate(workOnTaskData); err != nil {
		return c.JSON(http.StatusBadRequest, errInvalidTaskData)
	}

	err = h.service.WorkOnTask(workOnTaskData, userData.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err))
	}

	return c.JSON(http.StatusOK, true)
}
