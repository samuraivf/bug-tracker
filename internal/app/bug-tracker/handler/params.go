package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

//go:generate mockgen -source=params.go -destination=mocks/params.go

type Params interface {
	GetIdParam(c echo.Context) (uint64, error)
	GetUsernameParam(c echo.Context) (string, error)
}

type params struct{}

func (p *params) GetIdParam(c echo.Context) (uint64, error) {
	id := c.Param("id")

	if id == "" {
		return 0, errInvalidParam
	}

	uint64ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, errInvalidParam
	}

	return uint64ID, nil
}

func (p *params) GetUsernameParam(c echo.Context) (string, error) {
	username := c.Param("username")

	if username == "" {
		return "", errInvalidParam
	}

	return username, nil
}
