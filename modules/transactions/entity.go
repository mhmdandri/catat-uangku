package transactions

import (
	"catatan-keuangan/modules/attachments"
	transactionlines "catatan-keuangan/modules/transaction_lines"
	"time"
)

type Transactions struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID         *int      `gorm:"index" json:"group_id"`
	CategoryID      int       `gorm:"not null;index" json:"category_id"`
	CreatedByUserID int       `gorm:"not null;index" json:"created_by_user_id"`
	Date            time.Time `gorm:"not null" json:"date"`
	Type            string    `gorm:"type:varchar(50);not null" json:"type"`
	TotalAmount     float64   `gorm:"not null" json:"total_amount"`
	Scope           string    `gorm:"type:varchar(50);not null" json:"scope"`
	PrivateToUserID *int      `gorm:"index" json:"private_to_user_id,omitempty"`
	Description     *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	TransactionLines []transactionlines.TransactionLine `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"transaction_lines,omitempty"`
	Attachments      []attachments.Attachment           `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"attachments,omitempty"`
}
