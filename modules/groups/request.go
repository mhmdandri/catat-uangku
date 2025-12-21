package groups

import "github.com/google/uuid"

type GroupRequest struct {
	CreatorUserID uuid.UUID `json:"creator_user_id" swaggerignore:"true"`
	Name          string    `gorm:"type:varchar(100)" binding:"required" json:"name"`
	Type          string    `gorm:"type:varchar(50)" binding:"required" json:"type"`
}
