package transactions

import (
	"catatan-keuangan/modules/attachments"
	"catatan-keuangan/modules/categories"
	transactionlines "catatan-keuangan/modules/transaction_lines"

	"github.com/google/uuid"
)

type TransactionResponse struct {
	ID               uuid.UUID                                  `json:"id"`
	GroupID          *uuid.UUID                                 `json:"group_id"`
	CategoryID       uuid.UUID                                  `json:"category_id"`
	CreatedByUserID  uuid.UUID                                  `json:"created_by_user_id"`
	Title            string                                     `json:"title"`
	Date             string                                     `json:"date"`
	Type             string                                     `json:"type"`
	TotalAmount      float64                                    `json:"total_amount"`
	Scope            string                                     `json:"scope"`
	Description      *string                                    `json:"description,omitempty"`
	Category         categories.CategoryMini                    `json:"category"`
	TransactionLines []transactionlines.TransactionLineResponse `json:"transaction_lines,omitempty"`
	Attachments      []attachments.AttachmentResponse           `json:"attachments,omitempty"`
}

func FormatTransactionResponse(t Transactions) TransactionResponse {
	return TransactionResponse{
		ID:               t.ID,
		GroupID:          t.GroupID,
		CategoryID:       t.CategoryID,
		CreatedByUserID:  t.CreatedByUserID,
		Title:            t.Title,
		Date:             t.Date.Format("2006-01-02"),
		Type:             t.Type,
		TotalAmount:      t.TotalAmount,
		Scope:            t.Scope,
		Description:      t.Description,
		Category:         categories.FormatCategoryMini(t.Category),
		TransactionLines: transactionlines.FormatTransactionLineResponses(t.TransactionLines),
		Attachments:      attachments.FormatAttachmentResponses(t.Attachments),
	}
}

func FormatTransactionResponses(ts []Transactions) []TransactionResponse {
	formatted := make([]TransactionResponse, 0, len(ts))
	for _, t := range ts {
		formatted = append(formatted, FormatTransactionResponse(t))
	}
	return formatted
}
