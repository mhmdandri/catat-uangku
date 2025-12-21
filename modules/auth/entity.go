package auth

import (
	"catatan-keuangan/modules/users"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Token       string    `gorm:"type:varchar(255);not null"`
	ExpiresAt   time.Time `gorm:"not null"`
	RevokedAt   *time.Time
	RotatedFrom *uuid.UUID
	CreatedAt   time.Time

	User users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
