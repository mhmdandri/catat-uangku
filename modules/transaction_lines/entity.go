package transactionlines

type TransactionLine struct {
	ID            int     `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID int     `gorm:"not null;index" json:"transaction_id"`
	AccountID     int     `gorm:"not null;index" json:"account_id"`
	Debit         float64 `gorm:"not null;default:0" json:"debit"`
	Credit        float64 `gorm:"not null;default:0" json:"credit"`
	Note          *string `gorm:"type:text" json:"note,omitempty"`
}
