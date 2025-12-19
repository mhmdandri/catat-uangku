package categories

import "github.com/google/uuid"

type CategoryResponse struct {
	ID          uuid.UUID  `json:"id"`
	GroupID     *uuid.UUID `json:"group_id"`
	OwnerUserID *uuid.UUID `json:"owner_user_id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Color       *string    `json:"color,omitempty"`
	Icon        string     `json:"icon"`
}

type CategoryMini struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Type  string    `json:"type"`
	Color *string   `json:"color,omitempty"`
	Icon  string    `json:"icon"`
}

func FormatCategoryMini(category Category) CategoryMini {
	return CategoryMini{
		ID:    category.ID,
		Name:  category.Name,
		Type:  category.Type,
		Color: category.Color,
		Icon:  category.Icon,
	}
}

func FormatCategoryResponse(category Category) CategoryResponse {
	return CategoryResponse{
		ID:          category.ID,
		GroupID:     category.GroupID,
		OwnerUserID: category.OwnerUserID,
		Name:        category.Name,
		Type:        category.Type,
		Color:       category.Color,
		Icon:        category.Icon,
	}
}

func FormatCategoryResponses(categories []Category) []CategoryResponse {
	formatted := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		formatted = append(formatted, FormatCategoryResponse(category))
	}
	return formatted
}
