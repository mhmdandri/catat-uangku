package categories

type CategoryResponse struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	OwnerUserID int    `json:"owner_user_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

func FormatCategoryResponse(category Category) CategoryResponse {
	var groupID, ownerID int
	if category.GroupID != nil {
		groupID = *category.GroupID
	}
	if category.OwnerUserID != nil {
		ownerID = *category.OwnerUserID
	}

	return CategoryResponse{
		ID:          category.ID,
		GroupID:     groupID,
		OwnerUserID: ownerID,
		Name:        category.Name,
		Type:        category.Type,
	}
}

func FormatCategoryResponses(categories []Category) []CategoryResponse {
	formatted := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		formatted = append(formatted, FormatCategoryResponse(category))
	}
	return formatted
}
