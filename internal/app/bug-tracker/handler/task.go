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
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidJSON))
	}

	if err := c.Validate(workOnTaskData); err != nil {
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidTaskData))
	}

	err = h.service.WorkOnTask(workOnTaskData, userData.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err))
	}

	return c.JSON(http.StatusOK, true)
}

func (h *Handler) stopWorkOnTask(c echo.Context) error {
	workOnTaskData := new(dto.WorkOnTaskDto)
	userData, err := getUserData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	if err := c.Bind(workOnTaskData); err != nil {
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidJSON))
	}

	if err := c.Validate(workOnTaskData); err != nil {
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidTaskData))
	}

	err = h.service.StopWorkOnTask(workOnTaskData, userData.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err))
	}

	return c.JSON(http.StatusOK, true)
}

func (h *Handler) updateTask(c echo.Context) error {
	taskData := new(dto.UpdateTaskDto)
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

	id, err := h.service.Task.UpdateTask(taskData, userData.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err))
	}

	return c.JSON(http.StatusOK, id)
}

func (h *Handler) getTaskById(c echo.Context) error {
	id, err := h.params.GetIdParam(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	task, err := h.service.Task.GetTaskById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, newErrorMessage(errTaskNotFound))
	}

	return c.JSON(http.StatusFound, task)
}

func (h *Handler) getTaskByIdWithAssignee(c echo.Context) error {
	id, err := h.params.GetIdParam(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	task, err := h.service.Task.GetTaskById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, newErrorMessage(errTaskNotFound))
	}

	if !task.Assignee.Valid {
		return c.JSON(http.StatusFound, dto.TaskByIdWithAssignee{
			Task: task,
		})
	}

	assignee, err := h.service.GetUserById(uint64(task.Assignee.Int64))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(errTaskNotFound))
	}

	return c.JSON(http.StatusFound, dto.TaskByIdWithAssignee{
		Task:     task,
		Assignee: assignee,
	})
}

func (h *Handler) deleteTask(c echo.Context) error {
	taskData := new(dto.DeleteTaskDto)
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

	if err := h.service.Task.DeleteTask(taskData, userData.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err))
	}

	return c.JSON(http.StatusOK, true)
}
