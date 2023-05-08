package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
)

func (h *Handler) createProject(c echo.Context) error {
	projectData := new(dto.CreateProjectDto)

	if err := c.Bind(projectData); err != nil {
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidJSON))
	}

	if err := c.Validate(projectData); err != nil {
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidCreateProjectData))
	}

	id, err := h.service.Project.CreateProject(projectData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(errInternalServerError))
	}

	return c.JSON(http.StatusOK, id)
}

func (h *Handler) getProjectById(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errProjectNotFound))
	}

	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errProjectNotFound))
	}

	project, err := h.service.Project.GetProjectById(uint64ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, newErrorMessage(errProjectNotFound))
	}

	return c.JSON(http.StatusFound, project)
}
