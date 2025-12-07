package accounts

type AccountRequest struct {
	OwnerUserID int    `json:"owner_user_id" binding:"required"`
	GroupID     int    `json:"group_id"`
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Currency    string `json:"currency" binding:"required"`
	Scope       string `json:"scope" binding:"required"`
	IsShared    bool   `json:"is_shared"`
	IsActive    bool   `json:"is_active"`
}
