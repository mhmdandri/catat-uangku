package transactionlines

import "github.com/google/uuid"

type TransactionLineResponse struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	AccountID     uuid.UUID `json:"account_id"`
	AccountName   string    `json:"account_name"`
	Debit         float64   `json:"debit"`
	Credit        float64   `json:"credit"`
	Note          *string   `json:"note,omitempty"`
}

func FormatTransactionLineResponse(t TransactionLine) TransactionLineResponse {
	return TransactionLineResponse{
		ID:            t.ID,
		TransactionID: t.TransactionID,
		AccountID:     t.AccountID,
		AccountName:   t.Accounts.Name,
		Debit:         t.Debit,
		Credit:        t.Credit,
		Note:          t.Note,
	}
}

func FormatTransactionLineResponses(lines []TransactionLine) []TransactionLineResponse {
	formatted := make([]TransactionLineResponse, 0, len(lines))
	for _, line := range lines {
		formatted = append(formatted, FormatTransactionLineResponse(line))
	}
	return formatted
}
