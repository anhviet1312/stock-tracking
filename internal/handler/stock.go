package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"

	"codebase/internal/service/stock"
	httpx "codebase/pkg/httpx-echo"
)

type GroupStock struct {
	cfg *Config
}

func (group *GroupStock) GetStocksByGroupHandler(c echo.Context) error {
	ctx := c.Request().Context()
	
	groupParam := c.Param("group")
	if groupParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "group is required"})
	}

	srv, err := do.Invoke[stock.ServiceStock](group.cfg.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	items, err := srv.GetStockGroup(ctx, groupParam)
	return httpx.RestAbort(c, items, err)
}
