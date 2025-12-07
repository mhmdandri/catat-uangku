package categories

import (
	"catatan-keuangan/modules/transactions"
	"time"
)

type Category struct {
	ID           int                         `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID      *int                        `gorm:"index" json:"group_id"`
	OwnerUserID  *int                        `gorm:"index" json:"owner_user_id"`
	Name         string                      `gorm:"type:varchar(100);not null" json:"name"`
	Type         string                      `gorm:"type:varchar(50);not null" json:"type"`
	CreatedAt    time.Time                   `json:"created_at"`
	UpdatedAt    time.Time                   `json:"updated_at"`
	Transactions []transactions.Transactions `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"transactions,omitempty"`
}
