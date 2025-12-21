package accounts

import "github.com/google/uuid"

type AccountRequest struct {
	OwnerUserID  uuid.UUID  `json:"owner_user_id" swaggerignore:"true"`
	GroupID      *uuid.UUID `json:"group_id"`
	Name         string     `json:"name" binding:"required"`
	Type         string     `json:"type" binding:"required"`
	Number       *string    `json:"number"`
	FirstBalance *float64   `json:"first_balance" binding:"required,gte=0"`
	Currency     string     `json:"currency" binding:"required"`
	Scope        string     `json:"scope" binding:"required"`
	IsShared     bool       `json:"is_shared"`
	IsActive     bool       `json:"is_active"`
}

type AccountUpdateRequest struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Number   *string `json:"number"`
	IsShared bool    `json:"is_shared"`
	IsActive bool    `json:"is_active"`
}
