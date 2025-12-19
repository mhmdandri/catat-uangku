package categories

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	GroupID     *uuid.UUID `gorm:"type:uuid;index" json:"group_id"`
	OwnerUserID *uuid.UUID `gorm:"type:uuid;index" json:"owner_user_id"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Type        string     `gorm:"type:varchar(50);not null" json:"type"`
	Color       *string    `gorm:"type:varchar(7)" json:"color,omitempty"`
	Icon        string     `gorm:"type:varchar(50);default:'Circle'" json:"icon"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
