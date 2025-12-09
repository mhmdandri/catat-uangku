package userprofile

import (
	"time"

	"github.com/google/uuid"
)

type UserProfileResponse struct {
	ID         int64      `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	FirstName  string     `json:"first_name"`
	LastName   *string    `json:"last_name,omitempty"`
	Email      string     `json:"email"`
	Phone      *string    `json:"phone,omitempty"`
	Address    *string    `json:"address,omitempty"`
	Bio        *string    `json:"bio,omitempty"`
	Birthdate  *time.Time `json:"birthdate,omitempty"`
	Age        *int       `json:"age,omitempty"`
	IsVerified bool       `json:"is_verified"`
	AvatarURL  *string    `json:"avatar_url,omitempty"`
}

func FormatUserProfileResponse(profile UserProfile) UserProfileResponse {
	return UserProfileResponse{
		ID:         profile.ID,
		UserID:     profile.UserID,
		FirstName:  profile.FirstName,
		LastName:   profile.LastName,
		Email:      profile.Email,
		Phone:      profile.Phone,
		Address:    profile.Address,
		Bio:        profile.Bio,
		Birthdate:  profile.Birthdate,
		Age:        profile.Age,
		IsVerified: profile.IsVerified,
		AvatarURL:  profile.AvatarURL,
	}
}
