package handler

import (
	"github.com/labstack/echo/v4"
)

func setRoutes(e *echo.Echo, h *Handler) *echo.Echo {
	auth := e.Group(auth)
	{
		auth.POST(signUp, h.signUp)
		auth.POST(signIn, func(c echo.Context) error {
			return h.signIn(c, h.createTokens)
		}, h.isUnauthorized)
		auth.GET(refresh, func(c echo.Context) error {
			return h.refresh(c, h.createTokens)
		})
		auth.GET(logout, h.logout)
		auth.POST(verify, h.verifyEmail)
		auth.POST(setEmail, h.setEmail)
	}

	project := e.Group(project, h.isAuthorized)
	{
		project.POST(create, h.createProject)
		project.GET(id, h.getProjectById)
		project.DELETE(id, h.deleteProject)
		project.PUT(update, h.updateProject)
		project.POST(addMember, h.addMember)
		project.DELETE(deleteMember, h.deleteMember)
		project.GET(leave, h.leaveProject)
		project.POST(setAdmin, h.setNewAdmin)
	}

	task := e.Group(task, h.isAuthorized)
	{
		task.POST(create, h.createTask)
		task.POST(workOnTask, h.workOnTask)
	}

	return e
}
