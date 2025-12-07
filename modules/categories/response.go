package categories

import "github.com/google/uuid"

type CategoryResponse struct {
	ID          uuid.UUID  `json:"id"`
	GroupID     *uuid.UUID `json:"group_id"`
	OwnerUserID *uuid.UUID `json:"owner_user_id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
}

func FormatCategoryResponse(category Category) CategoryResponse {
	return CategoryResponse{
		ID:          category.ID,
		GroupID:     category.GroupID,
		OwnerUserID: category.OwnerUserID,
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
