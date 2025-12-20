package transactionlines

import (
	"catatan-keuangan/modules/accounts"

	"github.com/google/uuid"
)

type TransactionLine struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null;index" json:"transaction_id"`
	AccountID     uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	Debit         float64   `gorm:"not null;default:0" json:"debit"`
	Credit        float64   `gorm:"not null;default:0" json:"credit"`
	Note          *string   `gorm:"type:text" json:"note,omitempty"`

	Accounts accounts.Account `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"account,omitempty"`
}
