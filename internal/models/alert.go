package models

import (
	"time"

	"github.com/google/uuid"

	"codebase/internal/bob"
)

type UserStockAlert struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Symbol         string    `json:"symbol"`
	HighThreshold  *float64  `json:"high_threshold"`
	LowThreshold   *float64  `json:"low_threshold"`
	LastNotifiedAt *time.Time `json:"last_notified_at"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func UserStockAlertBobToRaw(item *bob.UserStockAlert) *UserStockAlert {
	if item == nil {
		return nil
	}
	var high *float64
	if v, ok := item.HighThreshold.Get(); ok {
		val, _ := v.Float64()
		high = &val
	}
	var low *float64
	if v, ok := item.LowThreshold.Get(); ok {
		val, _ := v.Float64()
		low = &val
	}
	var lastNotified *time.Time
	if v, ok := item.LastNotifiedAt.Get(); ok {
		lastNotified = &v
	}
	return &UserStockAlert{
		ID:             item.ID,
		UserID:         item.UserID,
		Symbol:         item.Symbol,
		HighThreshold:  high,
		LowThreshold:   low,
		LastNotifiedAt: lastNotified,
		IsActive:       item.IsActive,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func SliceUserStockAlertBobToRaw(items []*bob.UserStockAlert) []*UserStockAlert {
	result := make([]*UserStockAlert, 0, len(items))
	for _, item := range items {
		result = append(result, UserStockAlertBobToRaw(item))
	}
	return result
}

type SetStockAlertRequest struct {
	Symbol        string   `json:"symbol" validate:"required"`
	HighThreshold *float64 `json:"high_threshold"`
	LowThreshold  *float64 `json:"low_threshold"`
}

type UpdateTelegramChatIDRequest struct {
	TelegramChatID string `json:"telegram_chat_id" validate:"required"`
}
