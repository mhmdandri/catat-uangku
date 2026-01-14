package accounts

import (
	"sort"

	"github.com/google/uuid"
)

type AccountResponse struct {
	ID           uuid.UUID  `json:"id"`
	OwnerUserID  uuid.UUID  `json:"owner_user_id"`
	GroupID      *uuid.UUID `json:"group_id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Number       *string    `json:"number"`
	FirstBalance float64    `json:"first_balance"`
	Balance      float64    `json:"balance"`
	Currency     string     `json:"currency"`
	Scope        string     `json:"scope"`
	IsShared     bool       `json:"is_shared"`
	IsActive     bool       `json:"is_active"`
}

type AccountListItemResponse struct {
	ID       uuid.UUID  `json:"id"`
	GroupID  *uuid.UUID `json:"group_id"`
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	Number   *string    `json:"number"`
	Balance  float64    `json:"balance"`
	Currency string     `json:"currency"`
	Scope    string     `json:"scope"`
	IsShared bool       `json:"is_shared"`
	IsActive bool       `json:"is_active"`
}

type AccountSummaryResponse struct {
	Currency            string  `json:"currency"`
	TotalBalance        float64 `json:"totalBalance"`
	TotalBankBalance    float64 `json:"totalBankBalance"`
	TotalEWalletBalance float64 `json:"totalEWalletBalance"`
	TotalCashBalance    float64 `json:"totalCashBalance"`
}

func FormatAccountResponse(account Account) AccountResponse {
	return AccountResponse{
		ID:           account.ID,
		OwnerUserID:  account.OwnerUserID,
		GroupID:      account.GroupID,
		Name:         account.Name,
		Type:         account.Type,
		Number:       account.Number,
		FirstBalance: account.FirstBalance,
		Balance:      account.Balance,
		Currency:     account.Currency,
		Scope:        account.Scope,
		IsShared:     account.IsShared,
		IsActive:     account.IsActive,
	}
}

func FormatAccountListItemResponse(account Account) AccountListItemResponse {
	return AccountListItemResponse{
		ID:       account.ID,
		GroupID:  account.GroupID,
		Name:     account.Name,
		Type:     account.Type,
		Number:   account.Number,
		Balance:  account.Balance,
		Currency: account.Currency,
		Scope:    account.Scope,
		IsShared: account.IsShared,
		IsActive: account.IsActive,
	}
}

func FormatAccountListItemResponses(accounts []Account) []AccountListItemResponse {
	formatted := make([]AccountListItemResponse, 0, len(accounts))
	for _, account := range accounts {
		formatted = append(formatted, FormatAccountListItemResponse(account))
	}
	return formatted
}

func FormatAccountResponses(accounts []Account) []AccountResponse {
	formatted := make([]AccountResponse, 0, len(accounts))
	for _, account := range accounts {
		formatted = append(formatted, FormatAccountResponse(account))
	}
	return formatted
}

func BuildAccountSummaryByCurrency(accounts []Account) []AccountSummaryResponse {
	summaries := make(map[string]*AccountSummaryResponse)
	for _, account := range accounts {
		currency := account.Currency
		if currency == "" {
			currency = "IDR"
		}
		summary, ok := summaries[currency]
		if !ok {
			summary = &AccountSummaryResponse{Currency: currency}
			summaries[currency] = summary
		}
		balance := account.Balance
		summary.TotalBalance += balance
		switch account.Type {
		case "bank":
			summary.TotalBankBalance += balance
		case "e-wallet":
			summary.TotalEWalletBalance += balance
		case "cash":
			summary.TotalCashBalance += balance
		}
	}

	result := make([]AccountSummaryResponse, 0, len(summaries))
	for _, summary := range summaries {
		result = append(result, *summary)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Currency == "IDR" && result[j].Currency != "IDR" {
			return true
		}
		if result[j].Currency == "IDR" && result[i].Currency != "IDR" {
			return false
		}
		return result[i].Currency < result[j].Currency
	})
	return result
}
