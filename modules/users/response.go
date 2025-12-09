package users

import (
	"catatan-keuangan/modules/accounts"
	groupmembers "catatan-keuangan/modules/group_members"
	userprofile "catatan-keuangan/modules/user_profile"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID           uuid.UUID                          `json:"id"`
	Name         string                             `json:"name"`
	Email        string                             `json:"email"`
	CreatedAt    string                             `json:"created_at"`
	Accounts     []accounts.AccountResponse         `json:"accounts,omitempty"`
	GroupMembers []groupmembers.GroupMemberResponse `json:"group_members,omitempty"`
	Profile      *userprofile.UserProfileResponse   `json:"profile,omitempty"`
}

func FormatUserResponse(user User) UserResponse {
	var profile *userprofile.UserProfileResponse
	if user.Profile.ID != 0 {
		formatted := userprofile.FormatUserProfileResponse(user.Profile)
		profile = &formatted
	}
	return UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt.Format("2006-01-02 15:04:05"),
		Accounts:     accounts.FormatAccountResponses(user.Accounts),
		GroupMembers: groupmembers.FormatGroupMemberResponses(user.GroupMembers),
		Profile:      profile,
	}
}

func FormatUserResponses(users []User) []UserResponse {
	formatted := make([]UserResponse, 0, len(users))
	for _, user := range users {
		formatted = append(formatted, FormatUserResponse(user))
	}
	return formatted
}
