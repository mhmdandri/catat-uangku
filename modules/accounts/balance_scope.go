package accounts

import "gorm.io/gorm"

func WithBalance(tx *gorm.DB) *gorm.DB {
	return tx.Select(`
	accounts.*,
	COALESCE(accounts.first_balance,0) + COALESCE(SUM(tl.credit - tl.debit),0) AS balance`).
		Joins("LEFT JOIN transaction_lines tl ON tl.account_id = accounts.id").
		Group("accounts.id")
}
