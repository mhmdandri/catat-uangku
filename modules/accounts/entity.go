package accounts

import (
	transactionlines "catatan-keuangan/modules/transaction_lines"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OwnerUserID uuid.UUID  `gorm:"type:uuid;not null;index"`
	GroupID     *uuid.UUID `gorm:"type:uuid;index"`
	Name        string     `gorm:"type:varchar(100);not null"`
	Type        string     `gorm:"type:varchar(50);not null"`
	Currency    string     `gorm:"type:varchar(10);not null"`
	Scope       string     `gorm:"type:varchar(50);not null"`
	IsShared    bool       `gorm:"default:false"`
	IsActive    bool       `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	TransactionLines []transactionlines.TransactionLine `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"transaction_lines,omitempty"`
}
