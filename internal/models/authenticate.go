package models

import "github.com/google/uuid"

type UserClaim struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
}
