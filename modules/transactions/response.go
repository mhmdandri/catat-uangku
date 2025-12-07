package transactions

import (
	"catatan-keuangan/modules/attachments"
	transactionlines "catatan-keuangan/modules/transaction_lines"
)

type TransactionResponse struct {
	ID               int                                `json:"id"`
	GroupID          *int                               `json:"group_id"`
	CategoryID       int                                `json:"category_id"`
	CreatedByUserID  int                                `json:"created_by_user_id"`
	Date             string                             `json:"date"`
	Type             string                             `json:"type"`
	TotalAmount      float64                            `json:"total_amount"`
	Scope            string                             `json:"scope"`
	Description      *string                            `json:"description,omitempty"`
	TransactionLines []transactionlines.TransactionLine `json:"transaction_lines,omitempty"`
	Attachments      []attachments.AttachmentResponse   `json:"attachments,omitempty"`
}

func FormatTransactionResponse(t Transactions) TransactionResponse {
	return TransactionResponse{
		ID:               t.ID,
		GroupID:          t.GroupID,
		CategoryID:       t.CategoryID,
		CreatedByUserID:  t.CreatedByUserID,
		Date:             t.Date.Format("2006-01-02"),
		Type:             t.Type,
		TotalAmount:      t.TotalAmount,
		Scope:            t.Scope,
		Description:      t.Description,
		TransactionLines: t.TransactionLines,
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
