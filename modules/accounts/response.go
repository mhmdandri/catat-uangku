package accounts

import (
	transactionlines "catatan-keuangan/modules/transaction_lines"

	"github.com/google/uuid"
)

type AccountResponse struct {
	ID               uuid.UUID                          `json:"id"`
	OwnerUserID      uuid.UUID                          `json:"owner_user_id"`
	GroupID          *uuid.UUID                         `json:"group_id"`
	Name             string                             `json:"name"`
	Type             string                             `json:"type"`
	Number           *string                            `json:"number"`
	FirstBalance     float64                            `json:"first_balance"`
	Balance          float64                            `json:"balance"`
	Currency         string                             `json:"currency"`
	Scope            string                             `json:"scope"`
	IsShared         bool                               `json:"is_shared"`
	IsActive         bool                               `json:"is_active"`
	TransactionLines []transactionlines.TransactionLine `json:"transaction_lines,omitempty"`
}

func FormatAccountResponse(account Account) AccountResponse {
	return AccountResponse{
		ID:               account.ID,
		OwnerUserID:      account.OwnerUserID,
		GroupID:          account.GroupID,
		Name:             account.Name,
		Type:             account.Type,
		Number:           account.Number,
		FirstBalance:     account.FirstBalance,
		Balance:          account.Balance,
		Currency:         account.Currency,
		Scope:            account.Scope,
		IsShared:         account.IsShared,
		IsActive:         account.IsActive,
		TransactionLines: account.TransactionLines,
	}
}

func FormatAccountResponses(accounts []Account) []AccountResponse {
	formatted := make([]AccountResponse, 0, len(accounts))
	for _, account := range accounts {
		formatted = append(formatted, FormatAccountResponse(account))
	}
	return formatted
}
