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

func (group *GroupStock) GetStocksByExchangeHandler(c echo.Context) error {
	ctx := c.Request().Context()

	exchangeParam := c.Param("exchange")
	if exchangeParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "exchange is required"})
	}

	srv, err := do.Invoke[stock.ServiceStock](group.cfg.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	items, err := srv.GetStockExchange(ctx, exchangeParam)
	return httpx.RestAbort(c, items, err)
}

func (group *GroupStock) GetStockInfoHandler(c echo.Context) error {
	ctx := c.Request().Context()

	symbolParam := c.Param("symbol")
	if symbolParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "symbol is required"})
	}

	boardIDParam := c.QueryParam("boardId")

	srv, err := do.Invoke[stock.ServiceStock](group.cfg.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	item, err := srv.GetStockInfo(ctx, symbolParam, boardIDParam)
	return httpx.RestAbort(c, item, err)
}
