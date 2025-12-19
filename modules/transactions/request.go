package transactions

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type TransactionRequest struct {
	GroupID         *uuid.UUID `json:"group_id" form:"group_id"`
	CategoryID      uuid.UUID  `json:"category_id" form:"category_id" binding:"required"`
	CreatedByUserID uuid.UUID  `json:"created_by_user_id" form:"created_by_user_id" binding:"required"`
	AccountID       uuid.UUID  `json:"account_id" form:"account_id" binding:"required"`
	Title           string     `json:"title" form:"title" binding:"required"`
	Type            string     `json:"type" form:"type" binding:"required,oneof=income expense"`
	TotalAmount     float64    `json:"total_amount" form:"total_amount" binding:"required"`
	Scope           string     `json:"scope" form:"scope" binding:"required,oneof=personal group"`
	Description     *string    `json:"description" form:"description"`

	Attachments []*multipart.FileHeader `json:"attachments" form:"attachments"`
}
