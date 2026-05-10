package alert

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/samber/do"

	"codebase/internal/bob"
	"codebase/internal/models"
	"codebase/internal/service/utils"
	"github.com/shopspring/decimal"
)

type ServiceAlert interface {
	UpdateTelegramChatID(ctx context.Context, userID uuid.UUID, chatID string) error
	SetUserStockAlert(ctx context.Context, userID uuid.UUID, req models.SetStockAlertRequest) error
	RemoveUserStockAlert(ctx context.Context, userID uuid.UUID, alertID uuid.UUID) error
	ListUserStockAlerts(ctx context.Context, userID uuid.UUID) ([]*models.UserStockAlert, error)
	ListAllActiveStockAlerts(ctx context.Context) ([]*models.UserStockAlert, error)
}

type serviceAlert struct {
	Utils *utils.ServiceUtils
}

var _ ServiceAlert = (*serviceAlert)(nil)

func NewServiceAlert(container *do.Injector) (ServiceAlert, error) {
	utilsService, err := do.Invoke[*utils.ServiceUtils](container)
	if err != nil {
		return nil, err
	}

	return &serviceAlert{
		Utils: utilsService,
	}, nil
}

func (s *serviceAlert) UpdateTelegramChatID(ctx context.Context, userID uuid.UUID, chatID string) error {
	ds, err := s.Utils.Datastore.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = ds.Rollback(ctx) }()

	param := &bob.UserSetter{
		TelegramChatID: omitnull.From(chatID),
		UpdatedAt:      omit.From(time.Now()),
	}

	_, err = ds.UpdateUser(ctx, userID, param)
	if err != nil {
		return err
	}

	return ds.Commit(ctx)
}

func (s *serviceAlert) SetUserStockAlert(ctx context.Context, userID uuid.UUID, req models.SetStockAlertRequest) error {
	ds, err := s.Utils.Datastore.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = ds.Rollback(ctx) }()

	// Try to find if user already has an alert for this symbol
	alerts, err := ds.ListUserStockAlerts(ctx, userID)
	if err != nil {
		return err
	}

	var existingAlert *models.UserStockAlert
	for _, alert := range alerts {
		if alert.Symbol == req.Symbol {
			existingAlert = alert
			break
		}
	}

	if existingAlert != nil {
		var high, low omitnull.Val[decimal.Decimal]
		if req.HighThreshold != nil {
			high = omitnull.From(decimal.NewFromFloat(*req.HighThreshold))
		}
		if req.LowThreshold != nil {
			low = omitnull.From(decimal.NewFromFloat(*req.LowThreshold))
		}

		setter := &bob.UserStockAlertSetter{
			HighThreshold: high,
			LowThreshold:  low,
			UpdatedAt:     omit.From(time.Now()),
		}
		if err := ds.UpdateStockAlert(ctx, existingAlert.ID, setter); err != nil {
			return err
		}
	} else {
		var high, low omitnull.Val[decimal.Decimal]
		if req.HighThreshold != nil {
			high = omitnull.From(decimal.NewFromFloat(*req.HighThreshold))
		}
		if req.LowThreshold != nil {
			low = omitnull.From(decimal.NewFromFloat(*req.LowThreshold))
		}

		setter := &bob.UserStockAlertSetter{
			ID:            omit.From(uuid.New()),
			UserID:        omit.From(userID),
			Symbol:        omit.From(req.Symbol),
			HighThreshold: high,
			LowThreshold:  low,
			IsActive:      omit.From(true),
			CreatedAt:     omit.From(time.Now()),
			UpdatedAt:     omit.From(time.Now()),
		}
		if err := ds.SetStockAlert(ctx, setter); err != nil {
			return err
		}
	}

	return ds.Commit(ctx)
}

func (s *serviceAlert) RemoveUserStockAlert(ctx context.Context, userID uuid.UUID, alertID uuid.UUID) error {
	return s.Utils.Datastore.RemoveStockAlert(ctx, alertID, userID)
}

func (s *serviceAlert) ListUserStockAlerts(ctx context.Context, userID uuid.UUID) ([]*models.UserStockAlert, error) {
	return s.Utils.ReadOnlyDatastore.ListUserStockAlerts(ctx, userID)
}

func (s *serviceAlert) ListAllActiveStockAlerts(ctx context.Context) ([]*models.UserStockAlert, error) {
	return s.Utils.ReadOnlyDatastore.ListAllActiveStockAlerts(ctx)
}
