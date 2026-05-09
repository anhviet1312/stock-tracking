package handler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"

	"codebase/internal/models"
	"codebase/internal/service/alert"
	httpx "codebase/pkg/httpx-echo"
)

type GroupAlert struct {
	cfg *Config
}

func (h *GroupAlert) UpdateTelegramChatID(c echo.Context) error {
	userClaim, ok := c.Request().Context().Value(models.ContextUserClaimKey).(*models.UserClaim)
	if !ok || userClaim == nil {
		return httpx.RestAbort(c, nil, errors.New("unauthorized"))
	}

	var req models.UpdateTelegramChatIDRequest
	if err := c.Bind(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	if err := c.Validate(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	alertSvc := do.MustInvoke[alert.ServiceAlert](h.cfg.Container)
	if err := alertSvc.UpdateTelegramChatID(c.Request().Context(), userClaim.ID, req.TelegramChatID); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	return httpx.RestAbort(c, nil, nil)
}

func (h *GroupAlert) ListUserStockAlerts(c echo.Context) error {
	userClaim, ok := c.Request().Context().Value(models.ContextUserClaimKey).(*models.UserClaim)
	if !ok || userClaim == nil {
		return httpx.RestAbort(c, nil, errors.New("unauthorized"))
	}

	alertSvc := do.MustInvoke[alert.ServiceAlert](h.cfg.Container)
	alerts, err := alertSvc.ListUserStockAlerts(c.Request().Context(), userClaim.ID)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	return httpx.RestAbort(c, alerts, nil)
}

func (h *GroupAlert) SetUserStockAlert(c echo.Context) error {
	userClaim, ok := c.Request().Context().Value(models.ContextUserClaimKey).(*models.UserClaim)
	if !ok || userClaim == nil {
		return httpx.RestAbort(c, nil, errors.New("unauthorized"))
	}

	var req models.SetStockAlertRequest
	if err := c.Bind(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	if err := c.Validate(&req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	alertSvc := do.MustInvoke[alert.ServiceAlert](h.cfg.Container)
	if err := alertSvc.SetUserStockAlert(c.Request().Context(), userClaim.ID, req); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	return httpx.RestAbort(c, nil, nil)
}

func (h *GroupAlert) RemoveUserStockAlert(c echo.Context) error {
	userClaim, ok := c.Request().Context().Value(models.ContextUserClaimKey).(*models.UserClaim)
	if !ok || userClaim == nil {
		return httpx.RestAbort(c, nil, errors.New("unauthorized"))
	}

	alertIDStr := c.Param("id")
	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	alertSvc := do.MustInvoke[alert.ServiceAlert](h.cfg.Container)
	if err := alertSvc.RemoveUserStockAlert(c.Request().Context(), userClaim.ID, alertID); err != nil {
		return httpx.RestAbort(c, nil, err)
	}

	return httpx.RestAbort(c, nil, nil)
}
