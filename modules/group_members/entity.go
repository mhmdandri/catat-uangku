package groupmembers

import (
	"time"

	"github.com/google/uuid"
)

type GroupMember struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	GroupID  uuid.UUID `gorm:"type:uuid;not null;index" json:"group_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Role     string    `gorm:"type:varchar(50);not null" json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	IsActive bool      `gorm:"default:true" json:"is_active"`
}
