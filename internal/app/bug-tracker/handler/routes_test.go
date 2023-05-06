package handler

import (
	"sort"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Test_setRoutes(t *testing.T) {
	e1 := echo.New()
	e2 := echo.New()
	h := NewHandler(nil, nil, nil)

	auth := e1.Group(auth)
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

	project := e1.Group(project)
	{
		project.POST(create, h.createProject)
	}

	e2 = setRoutes(e2, h)

	e1Routes := e1.Routes()
	sort.SliceStable(e1Routes, func(i, j int) bool {
		return e1Routes[i].Name < e1Routes[j].Name
	})

	e2Routes := e2.Routes()
	sort.SliceStable(e2Routes, func(i, j int) bool {
		return e2Routes[i].Name < e2Routes[j].Name
	})

	require.Equal(t, len(e1Routes), len(e2Routes))

	for i, route := range e1Routes {
		require.Equal(t, route.Method, e2Routes[i].Method)
		require.Equal(t, route.Path, e2Routes[i].Path)
	}
}
