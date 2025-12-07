package groupmembers

import "time"

type GroupMember struct {
	ID       int       `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID  int       `gorm:"not null;index" json:"group_id"`
	UserID   int       `gorm:"not null;index" json:"user_id"`
	Role     string    `gorm:"type:varchar(50);not null" json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	IsActive bool      `gorm:"default:true" json:"is_active"`
}
