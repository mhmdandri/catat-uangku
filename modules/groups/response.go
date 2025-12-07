package groups

import groupmembers "catatan-keuangan/modules/group_members"

type GroupResponse struct {
	ID      int                                `json:"id"`
	Name    string                             `json:"name"`
	Type    string                             `json:"type"`
	Members []groupmembers.GroupMemberResponse `json:"members,omitempty"`
}

func FormatGroupResponse(group Group) GroupResponse {
	return GroupResponse{
		ID:      group.ID,
		Name:    group.Name,
		Type:    group.Type,
		Members: groupmembers.FormatGroupMemberResponses(group.GroupMembers),
	}
}

func FormatGroupResponses(groups []Group) []GroupResponse {
	formatted := make([]GroupResponse, 0, len(groups))
	for _, group := range groups {
		formatted = append(formatted, FormatGroupResponse(group))
	}
	return formatted
}
