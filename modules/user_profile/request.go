package userprofile

import "time"

type UpdateProfileRequest struct {
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	Phone     *string    `json:"phone"`
	Address   *string    `json:"address"`
	Bio       *string    `json:"bio"`
	Birthdate *time.Time `json:"birthdate"`
	Age       *int       `json:"age"`
	AvatarURL *string    `json:"avatar_url"`
}
