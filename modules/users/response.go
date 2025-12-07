package users

import (
	"catatan-keuangan/modules/accounts"
	groupmembers "catatan-keuangan/modules/group_members"
)

type UserResponse struct {
	ID           int                                `json:"id"`
	Name         string                             `json:"name"`
	Email        string                             `json:"email"`
	Accounts     []accounts.AccountResponse         `json:"accounts,omitempty"`
	GroupMembers []groupmembers.GroupMemberResponse `json:"group_members,omitempty"`
}

func FormatUserResponse(user User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Accounts:     accounts.FormatAccountResponses(user.Accounts),
		GroupMembers: groupmembers.FormatGroupMemberResponses(user.GroupMembers),
	}
}

func FormatUserResponses(users []User) []UserResponse {
	formatted := make([]UserResponse, 0, len(users))
	for _, user := range users {
		formatted = append(formatted, FormatUserResponse(user))
	}
	return formatted
}
