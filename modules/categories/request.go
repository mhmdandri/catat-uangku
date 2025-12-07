package categories

type CategoryRequest struct {
	GroupID     int    `json:"group_id" binding:"omitempty,gt=0"`
	OwnerUserID int    `json:"owner_user_id" binding:"omitempty,gt=0"`
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=income expense"`
}
