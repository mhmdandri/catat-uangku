package userprofile

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	FirstName  string     `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName   *string    `gorm:"type:varchar(100)" json:"last_name,omitempty"`
	Email      string     `gorm:"type:varchar(100);not null;uniqueIndex" json:"email"`
	Phone      *string    `gorm:"type:varchar(20);uniqueIndex" json:"phone,omitempty"`
	Address    *string    `gorm:"type:varchar(255)" json:"address,omitempty"`
	Bio        *string    `gorm:"type:text" json:"bio,omitempty"`
	Birthdate  *time.Time `json:"birthdate,omitempty"`
	Age        *int       `json:"age,omitempty"`
	IsVerified bool       `gorm:"default:false" json:"is_verified"`
	AvatarURL  *string    `gorm:"type:varchar(255)" json:"avatar_url,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
