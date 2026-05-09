package handler

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"

	"codebase/internal/service/stock"
	"codebase/pkg/errorx"
	httpx "codebase/pkg/httpx-echo"
)

type GroupStock struct {
	cfg *Config
}

func (group *GroupStock) GetStocksByGroupHandler(c echo.Context) error {
	ctx := c.Request().Context()

	groupParam := c.Param("group")
	if groupParam == "" {
		return httpx.RestAbort(c, nil, errorx.Wrap(errors.New("group is required"), errorx.Invalid))
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
		return httpx.RestAbort(c, nil, errorx.Wrap(errors.New("exchange is required"), errorx.Invalid))
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
		return httpx.RestAbort(c, nil, errorx.Wrap(errors.New("symbol is required"), errorx.Invalid))
	}

	boardIDParam := c.QueryParam("boardId")

	srv, err := do.Invoke[stock.ServiceStock](group.cfg.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	item, err := srv.GetStockInfo(ctx, symbolParam, boardIDParam)
	return httpx.RestAbort(c, item, err)
}

func (group *GroupStock) GetStockHistoryHandler(c echo.Context) error {
	ctx := c.Request().Context()

	symbolParam := c.Param("symbol")
	if symbolParam == "" {
		return httpx.RestAbort(c, nil, errorx.Wrap(errors.New("symbol is required"), errorx.Invalid))
	}

	resolutionParam := c.QueryParam("resolution")
	if resolutionParam == "" {
		resolutionParam = "1D"
	}

	fromParam, _ := strconv.ParseInt(c.QueryParam("from"), 10, 64)
	toParam, _ := strconv.ParseInt(c.QueryParam("to"), 10, 64)

	srv, err := do.Invoke[stock.ServiceStock](group.cfg.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	data, err := srv.GetStockHistory(ctx, symbolParam, resolutionParam, fromParam, toParam)
	return httpx.RestAbort(c, data, err)
}
