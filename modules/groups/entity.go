package groups

import (
	"catatan-keuangan/modules/accounts"
	"catatan-keuangan/modules/categories"
	groupmembers "catatan-keuangan/modules/group_members"
	"catatan-keuangan/modules/invitations"
	"catatan-keuangan/modules/transactions"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Type      string    `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	GroupMembers []groupmembers.GroupMember  `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"group_members,omitempty"`
	Invitations  []invitations.Invitation    `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"invitations,omitempty"`
	Accounts     []accounts.Account          `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"accounts,omitempty"`
	Categories   []categories.Category       `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"categories,omitempty"`
	Transactions []transactions.Transactions `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"transactions,omitempty"`
}
