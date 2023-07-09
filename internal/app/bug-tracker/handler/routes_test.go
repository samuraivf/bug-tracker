package handler

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Test_setRoutes(t *testing.T) {
	e := echo.New()
	expected := echo.New()
	h := &Handler{}

	auth := expected.Group(auth)
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

	project := expected.Group(project, h.isAuthorized)
	{
		project.POST(create, h.createProject)
		project.GET(id, h.getProjectById)
		project.GET(withTasks, h.getProjectByIdWithTasks)
		project.DELETE(id, h.deleteProject)
		project.PUT(update, h.updateProject)
		project.POST(addMember, h.addMember)
		project.DELETE(deleteMember, h.deleteMember)
		project.GET(leave, h.leaveProject)
		project.POST(setAdmin, h.setNewAdmin)
	}

	task := expected.Group(task, h.isAuthorized)
	{
		task.POST(create, h.createTask)
		task.POST(workOnTask, h.workOnTask)
		task.POST(stopWorkOnTask, h.stopWorkOnTask)
		task.PUT(update, h.updateTask)
		task.GET(id, h.getTaskById)
		task.DELETE(empty, h.deleteTask)
	}

	user := expected.Group(user, h.isAuthorized)
	{
		user.GET(id, h.getUserById)
	}

	e = setRoutes(e, h)

	require.Equal(t, len(expected.Routes()), len(e.Routes()))
}
