package handler

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"

	"codebase/internal/models"
	"codebase/internal/service/stock"
	"codebase/pkg/httpx-echo"
)

type GroupFavouriteStock struct {
	Config *Config
}

func (g *GroupFavouriteStock) AddFavouriteHandler(c echo.Context) error {
	userClaim, ok := c.Request().Context().Value(models.ContextUserClaimKey).(*models.UserClaim)
	if !ok || userClaim == nil {
		return httpx.RestAbort(c, nil, errors.New("unauthorized"))
	}

	var req models.AddFavouriteStockRequest
	if err := c.Bind(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	if err := c.Validate(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	svc, err := do.Invoke[stock.ServiceStock](g.Config.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	err = svc.AddFavouriteStock(c.Request().Context(), userClaim.ID, req.Symbol)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	return httpx.RestAbort(c, nil, nil)
}

func (g *GroupFavouriteStock) ListFavouritesHandler(c echo.Context) error {
	userClaim, ok := c.Request().Context().Value(models.ContextUserClaimKey).(*models.UserClaim)
	if !ok || userClaim == nil {
		return httpx.RestAbort(c, nil, errors.New("unauthorized"))
	}

	svc, err := do.Invoke[stock.ServiceStock](g.Config.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	res, err := svc.ListFavouriteStocks(c.Request().Context(), userClaim.ID)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	return httpx.RestAbort(c, res, nil)
}

func (g *GroupFavouriteStock) RemoveFavouriteHandler(c echo.Context) error {
	userClaim, ok := c.Request().Context().Value(models.ContextUserClaimKey).(*models.UserClaim)
	if !ok || userClaim == nil {
		return httpx.RestAbort(c, nil, errors.New("unauthorized"))
	}

	symbol := c.Param("symbol")
	if symbol == "" {
		return httpx.RestAbort(c, nil, errors.New("missing symbol parameter"))
	}

	svc, err := do.Invoke[stock.ServiceStock](g.Config.Container)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	err = svc.RemoveFavouriteStock(c.Request().Context(), userClaim.ID, symbol)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	return httpx.RestAbort(c, nil, nil)
}
