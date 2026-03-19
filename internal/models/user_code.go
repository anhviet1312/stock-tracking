package models

import (
	"time"

	"github.com/google/uuid"
)

type UserCodeType string

const (
	UserCodeTypeActivation UserCodeType = "activation"
)

type UserCodeStatus string

const (
	UserCodeStatusOrigin UserCodeStatus = "origin"
	UserCodeStatusUsed   UserCodeStatus = "used"
)

type UserCode struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.UUID      `json:"user_id"`
	Type        UserCodeType   `json:"type"`
	Status      UserCodeStatus `json:"status"`
	Value       string         `json:"-"`
	ExpiredTime time.Time      `json:"expired_time"`
	UpdatedAt   time.Time      `json:"-"`
	CreatedAt   time.Time      `json:"-"`
}
