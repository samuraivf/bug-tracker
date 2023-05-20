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
		project.DELETE(id, h.deleteProject)
		project.PUT(update, h.updateProject)
		project.POST(addMember, h.addMember)
		project.GET(leave, h.leaveProject)
	}

	e = setRoutes(e, h)

	require.Equal(t, len(expected.Routes()), len(e.Routes()))
}
