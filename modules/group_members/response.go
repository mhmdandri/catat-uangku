package groupmembers

type GroupMemberResponse struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	GroupID  int    `json:"group_id"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
	IsActive bool   `json:"is_active"`
}

func FormatGroupMemberResponse(gm GroupMember) GroupMemberResponse {
	return GroupMemberResponse{
		ID:       gm.ID,
		UserID:   gm.UserID,
		GroupID:  gm.GroupID,
		Role:     gm.Role,
		JoinedAt: gm.JoinedAt.Format("2006-01-02 15:04:05"),
		IsActive: gm.IsActive,
	}
}

func FormatGroupMemberResponses(gms []GroupMember) []GroupMemberResponse {
	formatted := make([]GroupMemberResponse, 0, len(gms))
	for _, gm := range gms {
		formatted = append(formatted, FormatGroupMemberResponse(gm))
	}
	return formatted
}
