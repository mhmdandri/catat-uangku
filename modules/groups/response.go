package groups

import (
	groupmembers "catatan-keuangan/modules/group_members"

	"github.com/google/uuid"
)

type GroupResponse struct {
	ID      uuid.UUID                          `json:"id"`
	Name    string                             `json:"name"`
	Type    string                             `json:"type"`
	Members []groupmembers.GroupMemberResponse `json:"members,omitempty"`
}

type GroupSummaryResponse struct {
	TotalGroups       int   `json:"totalGroups"`
	TotalMembers      int   `json:"totalMembers"`
	TotalTransactions int64 `json:"totalTransactions"`
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

func BuildGroupSummary(groups []Group, totalTransactions int64) GroupSummaryResponse {
	totalMembers := 0
	for _, group := range groups {
		for _, member := range group.GroupMembers {
			if member.IsActive {
				totalMembers++
			}
		}
	}
	return GroupSummaryResponse{
		TotalGroups:       len(groups),
		TotalMembers:      totalMembers,
		TotalTransactions: totalTransactions,
	}
}
