package accounts

type AccountResponse struct {
	ID          int    `json:"id"`
	OwnerUserID int    `json:"owner_user_id"`
	GroupID     int    `json:"group_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Currency    string `json:"currency"`
	Scope       string `json:"scope"`
	IsShared    bool   `json:"is_shared"`
	IsActive    bool   `json:"is_active"`
}

func FormatAccountResponse(account Account) AccountResponse {
	return AccountResponse{
		ID:          account.ID,
		OwnerUserID: account.OwnerUserID,
		GroupID:     account.GroupID,
		Name:        account.Name,
		Type:        account.Type,
		Currency:    account.Currency,
		Scope:       account.Scope,
		IsShared:    account.IsShared,
		IsActive:    account.IsActive,
	}
}

func FormatAccountResponses(accounts []Account) []AccountResponse {
	formatted := make([]AccountResponse, 0, len(accounts))
	for _, account := range accounts {
		formatted = append(formatted, FormatAccountResponse(account))
	}
	return formatted
}
