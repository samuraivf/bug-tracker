package handler

import (
	"net/http"

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
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidProjectData))
	}

	id, err := h.service.Project.CreateProject(projectData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(errInternalServerError))
	}

	return c.JSON(http.StatusOK, id)
}

func (h *Handler) getProjectById(c echo.Context) error {
	id, err := h.params.GetIdParam(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	project, err := h.service.Project.GetProjectById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, newErrorMessage(errProjectNotFound))
	}

	return c.JSON(http.StatusFound, project)
}

func (h *Handler) deleteProject(c echo.Context) error {
	id, err := h.params.GetIdParam(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	userData, err := getUserData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	if err := h.service.Project.DeleteProject(id, userData.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err))
	}

	return c.JSON(http.StatusOK, true)
}

func (h *Handler) updateProject(c echo.Context) error {
	projectData := new(dto.UpdateProjectDto)
	userData, err := getUserData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	if err := c.Bind(projectData); err != nil {
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidJSON))
	}

	err = h.service.Project.UpdateProject(projectData, userData.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(errInternalServerError))
	}

	return c.JSON(http.StatusOK, true)
}

func (h *Handler) addMember(c echo.Context) error {
	memberData := new(dto.AddMemberDto)
	userData, err := getUserData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	if err := c.Bind(memberData); err != nil {
		h.log.Error(err)
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidJSON))
	}

	if userData.UserID == memberData.MemberID {
		return c.JSON(http.StatusBadRequest, newErrorMessage(errInvalidOperation))
	}

	err = h.service.Project.AddMember(memberData, userData.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(errInternalServerError))
	}

	return c.JSON(http.StatusOK, true)
}

func (h *Handler) leaveProject(c echo.Context) error {
	id, err := h.params.GetIdParam(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	userData, err := getUserData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrorMessage(err))
	}

	err = h.service.Project.LeaveProject(id, userData.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorMessage(err))
	}

	return c.JSON(http.StatusOK, true)
}
