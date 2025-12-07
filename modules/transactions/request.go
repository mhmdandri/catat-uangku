package transactions

import "mime/multipart"

type TransactionRequest struct {
	GroupID         int     `json:"group_id" form:"group_id" binding:"omitempty,gt=0"`
	CategoryID      int     `json:"category_id" form:"category_id" binding:"required"`
	CreatedByUserID int     `json:"created_by_user_id" form:"created_by_user_id" binding:"required"`
	AccountID       int     `json:"account_id" form:"account_id" binding:"required"`
	Type            string  `json:"type" form:"type" binding:"required,oneof=income expense"`
	TotalAmount     float64 `json:"total_amount" form:"total_amount" binding:"required"`
	Scope           string  `json:"scope" form:"scope" binding:"required,oneof=personal group"`
	Description     *string `json:"description" form:"description"`

	Attachments []*multipart.FileHeader `json:"attachments" form:"attachments"`
}
