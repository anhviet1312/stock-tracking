package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"

	"codebase/internal/models"
	"codebase/internal/service/auth"
	"codebase/pkg/httpx-echo"
)

type GroupAuth struct {
	Config *Config
}

func (g *GroupAuth) RegisterHandler(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	if err := c.Validate(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	svc, err := do.Invoke[auth.ServiceAuth](g.Config.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	res, err := svc.Register(c.Request().Context(), &req)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	return httpx.RestAbort(c, res, nil)
}

func (g *GroupAuth) LoginHandler(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	if err := c.Validate(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	svc, err := do.Invoke[auth.ServiceAuth](g.Config.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	res, err := svc.Login(c.Request().Context(), &req)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	return httpx.RestAbort(c, res, nil)
}
