package categories

import "github.com/google/uuid"

type CategoryRequest struct {
	GroupID     *uuid.UUID `json:"group_id"`
	OwnerUserID *uuid.UUID `json:"owner_user_id"`
	Name        string     `json:"name" binding:"required"`
	Type        string     `json:"type" binding:"required,oneof=income expense"`
	Color       *string    `json:"color,omitempty" form:"color,omitempty"`
	Icon        string     `json:"icon" form:"icon" binding:"required"`
}
