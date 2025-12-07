package groups

type GroupRequest struct {
	CreatorUserID int    `json:"creator_user_id" binding:"required"`
	Name          string `gorm:"type:varchar(100)" binding:"required" json:"name"`
	Type          string `gorm:"type:varchar(50)" binding:"required" json:"type"`
}
