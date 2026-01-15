package auth

import "github.com/google/uuid"

type MeSummary struct {
	TotalTransaction int64 `json:"totalTransaction"`
	TotalAccount     int64 `json:"totalAccount"`
	TotalGroup       int64 `json:"totalGroup"`
	DurationMember   int64 `json:"durationMember"`
}

type MeUserProfile struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatarUrl"`
}

type MeData struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type MeResponse struct {
	Summary     MeSummary     `json:"summary"`
	UserProfile MeUserProfile `json:"userProfile"`
	Data        MeData        `json:"data"`
}
