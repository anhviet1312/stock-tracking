package models

import (
	"time"

	"github.com/google/uuid"

	"codebase/internal/bob"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserBobToRaw(item *bob.User) *User {
	if item == nil {
		return nil
	}
	return &User{
		ID:        item.ID,
		FirstName: item.FirstName.GetOrZero(),
		LastName:  item.LastName.GetOrZero(),
		Username:  item.Username,
		Email:     item.Email,
		IsActive:  item.IsActive,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

type UserFavouriteStock struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Symbol    string    `json:"symbol"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Optional relation field
	StockInfo *Stock `json:"stock_info,omitempty"`
}

func UserFavouriteStockBobToRaw(item *bob.UserFavouriteStock) *UserFavouriteStock {
	if item == nil {
		return nil
	}
	return &UserFavouriteStock{
		ID:        item.ID,
		UserID:    item.UserID,
		Symbol:    item.Symbol,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func SliceUserFavouriteStockBobToRaw(items []*bob.UserFavouriteStock) []*UserFavouriteStock {
	result := make([]*UserFavouriteStock, 0, len(items))
	for _, item := range items {
		result = append(result, UserFavouriteStockBobToRaw(item))
	}
	return result
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type AddFavouriteStockRequest struct {
	Symbol string `json:"symbol" validate:"required"`
}
