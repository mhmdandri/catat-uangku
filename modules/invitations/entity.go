package invitations

import "time"

type Invitation struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID       int       `gorm:"not null;index" json:"group_id"`
	InviterUserID int       `gorm:"not null;index" json:"inviter_user_id"`
	InviteeEmail  string    `gorm:"type:varchar(100);not null" json:"invitee_email"`
	Token         string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"token"`
	Role          string    `gorm:"type:varchar(50);not null" json:"role"`
	Status        string    `gorm:"type:varchar(50);not null" json:"status"`
	ExpiredAt     time.Time `json:"expired_at"`
	CreatedAt     time.Time `json:"created_at"`
}
